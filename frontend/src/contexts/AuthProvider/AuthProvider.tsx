import { createContext, ReactNode, useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { useLocalStorage } from '../../hooks/useLocalStorage';

export interface User {
  id: string;
  name?: string;
  photoUri?: string;
  role?: string;
}

export interface AuthContextType {
  accessToken: string;
  currentUser: User | null;
  logout: () => void;
  namespaceId: string;
  setAccessToken: (token: string) => void;
  setCurrentUser: (user: User | null) => void;
}

interface AuthContextProviderProps {
  children: ReactNode;
}

interface NamespaceIdResponse {
  namespaceId: string;
}

interface RefreshTokenResponse {
  accessToken: string;
  user: User;
}

export const AuthContext = createContext<AuthContextType | null>(null);

export const AuthContextProvider = ({ children }: AuthContextProviderProps) => {
  const [accessToken, setAccessToken] = useState('');
  const [currentUser, setCurrentUser] = useState<User | null>(
    localStorage.getItem('user')
      ? JSON.parse(localStorage.getItem('user') as string)
      : null
  );
  const [namespaceId, setNamespaceId] = useState<string>('');

  const { removeValue, setValue } = useLocalStorage('user', null);

  useEffect(() => {
    if (!accessToken && currentUser) refreshToken();

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (!accessToken) return;

    const getNamespaceId = async () => {
      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/namespace/search-by-name?name=default`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`
          }
        }
      );
      const data: NamespaceIdResponse = await response.json();

      return data;
    };

    const fetchNamespaceId = async () => {
      const { namespaceId } = await getNamespaceId();
      setNamespaceId(namespaceId);
    };

    fetchNamespaceId();
  }, [accessToken]);

  const logout = async () => {
    try {
      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/auth/logout`,
        {
          credentials: 'include',
          headers: {
            Authorization: `Bearer ${accessToken}`
          }
        }
      );

      if (!response.ok) throw new Error('Failed to logout');

      setCurrentUser(null);
      setAccessToken('');
      removeValue();
    } catch (error) {
      toast.error('Failed to logout');
    }
  };

  const refreshToken = async () => {
    try {
      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/auth/token/refresh`,
        {
          credentials: 'include'
        }
      );

      if (!response.ok) throw new Error('Failed to refresh token');

      const data: RefreshTokenResponse = await response.json();

      setCurrentUser(data.user);
      setAccessToken(data.accessToken);

      data.user && setValue(data.user);
    } catch (error) {
      setCurrentUser(null);
      removeValue();
    }
  };

  const VALUE = {
    accessToken,
    currentUser,
    logout,
    namespaceId,
    setAccessToken,
    setCurrentUser
  };

  return <AuthContext.Provider value={VALUE}>{children}</AuthContext.Provider>;
};
