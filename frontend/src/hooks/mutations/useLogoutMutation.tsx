import type {
  RefreshAccessTokenFn,
  User
} from '../../contexts/AuthProvider/AuthProvider';
import { fetchDataWithAuth } from '../utils';
import { toast } from 'react-toastify';
import { useLocalStorage } from '../useLocalStorage';
import { useMutation } from '@tanstack/react-query';

interface LogoutMutationProps {
  accessToken: string;
  refreshAccessToken: RefreshAccessTokenFn;
  setAccessToken: (accessToken: string) => void;
  setCurrentUser: (user: User | null) => void;
}

interface LogoutProps {
  accessToken: string;
  refreshAccessToken: RefreshAccessTokenFn;
}

const logout = ({ accessToken, refreshAccessToken }: LogoutProps) => {
  if (!accessToken) throw new Error('Something went wrong!');

  const url = `${import.meta.env.VITE_SERVER_API_URL}/auth/logout`;
  const options = {
    credentials: 'include' as RequestCredentials,
    headers: {
      Authorization: `Bearer ${accessToken}`
    }
  };

  return fetchDataWithAuth({ options, refreshAccessToken, url });
};

export const useLogoutMutation = ({
  accessToken,
  refreshAccessToken,
  setAccessToken,
  setCurrentUser
}: LogoutMutationProps) => {
  const { removeValue } = useLocalStorage('user', null);

  return useMutation({
    mutationFn: () => logout({ accessToken, refreshAccessToken }),
    onError: () => {
      toast.error('Failed to logout');
    },
    onSuccess: () => {
      setAccessToken('');
      setCurrentUser(null);
      removeValue();
    }
  });
};
