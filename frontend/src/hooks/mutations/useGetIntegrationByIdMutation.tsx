import type { Integration } from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import { fetchDataWithAuth } from '../utils';
import { toast } from 'react-toastify';
import { useAuthContext } from '../useAuthContext';
import { useEffect, useRef } from 'react';
import { useIntegrationsContext } from '../useIntegrationsContext';
import { useMutation } from '@tanstack/react-query';

interface GetIntegrationByIdProps {
  accessToken: string;
  integrationId: string;
  namespaceId: string;
  refreshAccessToken: RefreshAccessTokenFn;
  signal: AbortSignal;
}

interface UseGetIntegrationMutationProps {
  integrationId: string;
}

interface InstalledIntegrationResponse {
  integration: Integration;
}

const getIntegrationById = async ({
  accessToken,
  integrationId,
  namespaceId,
  refreshAccessToken,
  signal
}: GetIntegrationByIdProps) => {
  if (!accessToken || !namespaceId) throw new Error('Something went wrong!');

  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/integration/${integrationId}`;
  const options = {
    headers: {
      Authorization: `Bearer ${accessToken}`
    },
    signal
  };

  return fetchDataWithAuth<InstalledIntegrationResponse>({
    options,
    refreshAccessToken,
    url
  });
};

export const useGetIntegrationByIdMutation = ({
  integrationId
}: UseGetIntegrationMutationProps) => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();
  const { isModalOpen, setIsModalOpen, setSelectedIntegration } =
    useIntegrationsContext();

  const abortControllerRef = useRef<AbortController | null>(null);

  useEffect(() => {
    if (!isModalOpen && abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
  }, [isModalOpen]);

  return useMutation({
    mutationFn: () => {
      abortControllerRef.current = new AbortController();
      const { signal } = abortControllerRef.current;

      return getIntegrationById({
        accessToken,
        integrationId,
        namespaceId,
        refreshAccessToken,
        signal
      });
    },
    onError: error => {
      if (error.name === 'AbortError') return;

      toast.error('Failed to get integration data, please try again later');
      setIsModalOpen(false);
    },
    onMutate: () => {
      setIsModalOpen(true);
    },
    onSuccess: data => {
      isModalOpen && setSelectedIntegration(data.integration);
    }
  });
};
