import type { ConfigData, NewIntegrationPayload } from './types';
import type { Integration } from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import { fetchDataWithAuth } from '../utils';
import { toast } from 'react-toastify';
import { useAuthContext } from '../useAuthContext';
import { useIntegrationsContext } from '../useIntegrationsContext';
import { useMutation, useQueryClient } from '@tanstack/react-query';

interface CreateIntegrationProps {
  accessToken: string;
  namespaceId: string;
  newIntegration: NewIntegrationPayload;
  refreshAccessToken: RefreshAccessTokenFn;
}

interface InstallIntegrationResponse {
  configData: ConfigData | null;
  integration: Integration;
}

interface UseCreateIntegrationMutationProps {
  setConfigData: (configData: ConfigData | null) => void;
  setError: (error: Error | null) => void;
  setInstallationStep: (installationStep: 0 | 1) => void;
}

const createIntegration = ({
  accessToken,
  newIntegration,
  namespaceId,
  refreshAccessToken
}: CreateIntegrationProps) => {
  if (!accessToken || !namespaceId) throw new Error('Something went wrong!');

  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/integration`;
  const options = {
    body: JSON.stringify(newIntegration),
    headers: {
      Authorization: `Bearer ${accessToken}`,
      'Content-Type': 'application/json'
    },
    method: 'POST'
  };

  return fetchDataWithAuth<InstallIntegrationResponse>({
    options,
    refreshAccessToken,
    url
  });
};

export const useCreateIntegrationMutation = ({
  setConfigData,
  setError,
  setInstallationStep
}: UseCreateIntegrationMutationProps) => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();
  const { setIsModalOpen } = useIntegrationsContext();

  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (newIntegration: NewIntegrationPayload) =>
      createIntegration({
        accessToken,
        newIntegration,
        namespaceId,
        refreshAccessToken
      }),
    onError: error => {
      toast.error(`Failed to install integration`);

      if (error instanceof Error) {
        setError(error);
      } else {
        setError(new Error('An unknown error occurred'));
      }
    },
    onMutate: () => {
      setError(null);
    },
    onSuccess: data => {
      if (data.configData) {
        setConfigData(data.configData);
        setInstallationStep(1);
      } else {
        setInstallationStep(0);
        setIsModalOpen(false);
      }
      queryClient.invalidateQueries({ queryKey: ['installedIntegrations'] });

      toast.success(`Integration installed successfully`);
    }
  });
};
