import { createContext, ReactNode, useEffect, useState } from 'react';
import { useGetNamespaceMutation } from '../../hooks/mutations/useGetNamespaceMutation';
import { useLogoutMutation } from '../../hooks/mutations/useLogoutMutation';
import { useRefreshAccessTokenMutation } from '../../hooks/mutations/useRefreshTokenMutation';

export type RefreshAccessTokenFn = () => Promise<string>;

interface AuthContextProviderProps {
  children: ReactNode;
}
export interface AuthContextType {
  accessToken: string;
  currentUser: User | null;
  logout: () => void;
  namespaceId: string;
  refreshAccessToken: RefreshAccessTokenFn;
  setAccessToken: (token: string) => void;
  setCurrentUser: (user: User | null) => void;
}

export interface User {
  id: string;
  name?: string;
  photoUri?: string;
  role?: string;
}

export const AuthContext = createContext<AuthContextType | null>(null);

export const AuthContextProvider = ({ children }: AuthContextProviderProps) => {
  const [accessToken, setAccessToken] = useState('');
  const [currentUser, setCurrentUser] = useState<User | null>(
    localStorage.getItem('user')
      ? JSON.parse(localStorage.getItem('user') as string)
      : null
  );
  const [namespaceId, setNamespaceId] = useState('');

  const refreshAccessToken = async () => {
    try {
      const { accessToken } = await refreshAccessTokenMutateAsync();
      return accessToken;
    } catch (error) {
      console.error(error);
      return '';
    }
  };

  const { mutate: getNamespaceIdMutate } = useGetNamespaceMutation({
    accessToken,
    refreshAccessToken,
    setNamespaceId
  });

  const { isPending: isLogoutPending, mutate: logoutMutate } =
    useLogoutMutation({
      accessToken,
      refreshAccessToken,
      setAccessToken,
      setCurrentUser
    });

  const {
    mutate: refreshAccessTokenMutate,
    mutateAsync: refreshAccessTokenMutateAsync
  } = useRefreshAccessTokenMutation({
    setAccessToken,
    setCurrentUser
  });

  useEffect(() => {
    if (!accessToken && currentUser) refreshAccessTokenMutate();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (!accessToken || isLogoutPending) return;

    getNamespaceIdMutate();
  }, [accessToken, getNamespaceIdMutate, isLogoutPending]);

  const VALUE = {
    accessToken,
    currentUser,
    logout: logoutMutate,
    namespaceId,
    refreshAccessToken,
    setAccessToken,
    setCurrentUser
  };

  return <AuthContext.Provider value={VALUE}>{children}</AuthContext.Provider>;
};
