import { createContext, ReactNode, useEffect, useState } from 'react';

export interface User {
  id: string;
  name?: string;
  photoUri?: string;
  role?: string;
}

export interface AuthContextType {
  accessToken: string;
  currentUser: User | null;
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
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [namespaceId, setNamespaceId] = useState<string>('');

  useEffect(() => {
    if (import.meta.env.VITE_SKIP_AUTH === 'true') {
      setCurrentUser({
        id: '1',
        name: 'Default',
        role: 'user'
      });
    }

    const fetchNamespaceId = async () => {
      const { namespaceId } = await getNamespaceId();
      setNamespaceId(namespaceId);
    };

    const user = localStorage.getItem('user');

    if (user) {
      setCurrentUser(JSON.parse(user));
    }

    if (!accessToken && user) refreshToken();

    fetchNamespaceId();
  }, [accessToken]);

  const getNamespaceId = async () => {
    const response = await fetch(
      `${import.meta.env.VITE_SERVER_API_URL}/namespace/search-by-name?name=default`
    );
    const data: NamespaceIdResponse = await response.json();

    return data;
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

      data.user && localStorage.setItem('user', JSON.stringify(data.user));
    } catch (error) {
      setCurrentUser(null);
      localStorage.removeItem('user');
    }
  };

  const VALUE = {
    accessToken,
    currentUser,
    namespaceId,
    setAccessToken,
    setCurrentUser
  };

  return <AuthContext.Provider value={VALUE}>{children}</AuthContext.Provider>;
};
