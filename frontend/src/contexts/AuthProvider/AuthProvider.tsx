import { createContext, ReactNode, useEffect, useState } from 'react';

type UserType = 'github' | 'google' | 'signal0ne';

interface User {
  email: string;
  id: string;
  type: UserType;
  userName?: string;
}

export interface AuthContextType {
  currentUser: User | null;
  namespaceId: string;
  setCurrentUser: (user: User | null) => void;
}

interface AuthContextProviderProps {
  children: ReactNode;
}

interface NamespaceIdResponse {
  namespaceId: string;
}

export const AuthContext = createContext<AuthContextType | null>(null);

export const AuthContextProvider = ({ children }: AuthContextProviderProps) => {
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [namespaceId, setNamespaceId] = useState<string>('');

  const getNamespaceId = async () => {
    const response = await fetch(
      `${import.meta.env.VITE_SERVER_API_URL}/namespace/search-by-name?name=default`
    );
    const data: NamespaceIdResponse = await response.json();

    return data;
  };

  useEffect(() => {
    if (import.meta.env.VITE_SKIP_AUTH === 'true') {
      setCurrentUser({
        email: '',
        id: '1',
        type: 'signal0ne',
        userName: 'Default'
      });
    }

    const fetchNamespaceId = async () => {
      const { namespaceId } = await getNamespaceId();
      setNamespaceId(namespaceId);
    };

    fetchNamespaceId();
  }, []);

  const VALUE = {
    currentUser,
    namespaceId,
    setCurrentUser
  };

  return <AuthContext.Provider value={VALUE}>{children}</AuthContext.Provider>;
};
