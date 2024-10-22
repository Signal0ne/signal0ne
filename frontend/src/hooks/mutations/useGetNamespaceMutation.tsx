import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import { fetchDataWithAuth } from '../utils';
import { toast } from 'react-toastify';
import { useMutation } from '@tanstack/react-query';

interface GetNamespaceIdProps {
  accessToken: string;
  refreshAccessToken: RefreshAccessTokenFn;
}

interface GetNamespaceMutationProps {
  accessToken: string;
  refreshAccessToken: RefreshAccessTokenFn;
  setNamespaceId: (namespaceId: string) => void;
}

interface NamespaceIdResponse {
  namespaceId: string;
}

const getNamespaceId = ({
  accessToken,
  refreshAccessToken
}: GetNamespaceIdProps): Promise<NamespaceIdResponse> => {
  if (!accessToken) throw new Error('Something went wrong!');

  const url = `${import.meta.env.VITE_SERVER_API_URL}/namespace/search-by-name?name=default`;
  const options = {
    headers: {
      Authorization: `Bearer ${accessToken}`
    }
  };

  return fetchDataWithAuth<NamespaceIdResponse>({
    options,
    refreshAccessToken,
    url
  });
};

export const useGetNamespaceMutation = ({
  accessToken,
  refreshAccessToken,
  setNamespaceId
}: GetNamespaceMutationProps) =>
  useMutation({
    mutationFn: () => getNamespaceId({ accessToken, refreshAccessToken }),
    onError: () => {
      toast.error('Failed to get namespace ID');
      setNamespaceId('');
    },
    onSuccess: data => {
      setNamespaceId(data.namespaceId);
    }
  });
