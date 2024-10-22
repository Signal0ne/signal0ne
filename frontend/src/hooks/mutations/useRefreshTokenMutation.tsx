import type { User } from '../../contexts/AuthProvider/AuthProvider';
import { toast } from 'react-toastify';
import { useLocalStorage } from '../useLocalStorage';
import { useMutation } from '@tanstack/react-query';

interface RefreshAccessTokenMutationProps {
  setAccessToken: (accessToken: string) => void;
  setCurrentUser: (user: User | null) => void;
}

export interface RefreshAccessTokenResponse {
  accessToken: string;
  user: User;
}

const refreshAccessToken = async (): Promise<RefreshAccessTokenResponse> => {
  const url = `${import.meta.env.VITE_SERVER_API_URL}/auth/token/refresh`;
  const options = {
    credentials: 'include' as RequestCredentials
  };

  const response = await fetch(url, options);

  if (!response.ok) throw new Error('Failed to refresh token');

  const data: RefreshAccessTokenResponse = await response.json();

  return data;
};

export const useRefreshAccessTokenMutation = ({
  setAccessToken,
  setCurrentUser
}: RefreshAccessTokenMutationProps) => {
  const { removeValue, setValue } = useLocalStorage('user', null);

  return useMutation({
    mutationFn: refreshAccessToken,
    onError: () => {
      setAccessToken('');
      setCurrentUser(null);
      removeValue();

      if (toast.isActive('session-expired-toast')) return;

      toast.error('Session expired. Please login again', {
        toastId: 'session-expired-toast'
      });
    },
    onSuccess: data => {
      setAccessToken(data.accessToken);
      setCurrentUser(data.user);
      data.user && setValue(data.user);
    },
    retry: 3
  });
};
