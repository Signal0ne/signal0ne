import type { AuthPayload, AuthResponse } from './types';
import { toast } from 'react-toastify';
import { useAuthContext } from '../useAuthContext';
import { useLocalStorage } from '../useLocalStorage';
import { useMutation } from '@tanstack/react-query';

const login = async ({
  password,
  username
}: AuthPayload): Promise<AuthResponse> => {
  const url = `${import.meta.env.VITE_SERVER_API_URL}/auth/login`;
  const options = {
    body: JSON.stringify({ password, username }),
    credentials: 'include' as RequestCredentials,
    headers: {
      'Content-Type': 'application/json'
    },
    method: 'POST'
  };

  const response = await fetch(url, options);

  if (!response.ok) throw new Error('Failed to login');

  const data = await response.json();

  return data;
};

export const useLoginMutation = () => {
  const { setAccessToken, setCurrentUser } = useAuthContext();
  const { setValue } = useLocalStorage('user', null);

  return useMutation({
    mutationFn: login,
    onError: error => {
      if (error instanceof Error) {
        toast.error(error.message);
      } else {
        toast.error('An unknown error occurred. Please try again later.');
      }
    },
    onSuccess: data => {
      setAccessToken(data.accessToken);
      setCurrentUser(data.user);
      setValue(data.user);
    }
  });
};
