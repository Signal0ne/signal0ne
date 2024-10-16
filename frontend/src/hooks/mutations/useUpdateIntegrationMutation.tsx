import type { ConfigData, NewIntegrationPayload } from './types';
import type {
  InstalledIntegration,
  Integration
} from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import { fetchDataWithAuth } from '../utils';
import { toast } from 'react-toastify';
import { useAuthContext } from '../useAuthContext';
import { useIntegrationsContext } from '../useIntegrationsContext';
import { useMutation, useQueryClient } from '@tanstack/react-query';

interface InstallIntegrationResponse {
  configData: ConfigData | null;
  integration: Integration;
}

interface UpdateIntegrationProps {
  accessToken: string;
  namespaceId: string;
  newIntegration: NewIntegrationPayload;
  refreshAccessToken: RefreshAccessTokenFn;
  selectedIntegrationId: string;
}

interface UseUpdateIntegrationMutationProps {
  setConfigData: (configData: ConfigData | null) => void;
  setError: (error: Error | null) => void;
  setInstallationStep: (installationStep: 0 | 1) => void;
}

const updateIntegration = ({
  accessToken,
  newIntegration,
  namespaceId,
  refreshAccessToken,
  selectedIntegrationId
}: UpdateIntegrationProps) => {
  if (!accessToken || !namespaceId) throw new Error('Something went wrong!');

  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/integration/${selectedIntegrationId}`;
  const options = {
    body: JSON.stringify(newIntegration),
    headers: {
      Authorization: `Bearer ${accessToken}`,
      'Content-Type': 'application/json'
    },
    method: 'PATCH'
  };

  return fetchDataWithAuth<InstallIntegrationResponse>({
    options,
    refreshAccessToken,
    url
  });
};

export const useUpdateIntegrationMutation = ({
  setConfigData,
  setError,
  setInstallationStep
}: UseUpdateIntegrationMutationProps) => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();
  const { selectedIntegration, setIsModalOpen } = useIntegrationsContext();

  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (newIntegration: NewIntegrationPayload) =>
      updateIntegration({
        accessToken,
        newIntegration,
        namespaceId,
        refreshAccessToken,
        selectedIntegrationId: (selectedIntegration as InstalledIntegration).id
      }),
    onError: error => {
      toast.error(`Failed to update integration`);

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

      toast.success(`Integration updated successfully`);
    }
  });
};
