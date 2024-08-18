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
  setCurrentUser: (user: User | null) => void;
}

interface AuthContextProviderProps {
  children: ReactNode;
}

export const AuthContext = createContext<AuthContextType | null>(null);

export const AuthContextProvider = ({ children }: AuthContextProviderProps) => {
  const [currentUser, setCurrentUser] = useState<User | null>(null);

  useEffect(() => {
    if (import.meta.env.VITE_SKIP_AUTH === 'true') {
      setCurrentUser({
        email: '',
        id: '1',
        type: 'signal0ne',
        userName: 'Default'
      });
    }
  }, []);

  const VALUE = {
    currentUser,
    setCurrentUser
  };

  return <AuthContext.Provider value={VALUE}>{children}</AuthContext.Provider>;
};
