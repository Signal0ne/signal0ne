import type { Integration } from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import { fetchDataWithAuth } from '../utils';
import { useAuthContext } from '../useAuthContext';
import { useQuery } from '@tanstack/react-query';

interface FetchInstallableIntegrationsProps {
  accessToken: string;
  namespaceId: string;
  refreshAccessToken: RefreshAccessTokenFn;
  signal: AbortSignal;
}

interface FetchInstallableIntegrationsResponse {
  installableIntegrations: Integration[];
}

function fetchInstallableIntegrations({
  accessToken,
  namespaceId,
  refreshAccessToken,
  signal
}: FetchInstallableIntegrationsProps) {
  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/integration/installable`;
  const options = {
    headers: {
      Authorization: `Bearer ${accessToken}`
    },
    signal
  };

  return fetchDataWithAuth<FetchInstallableIntegrationsResponse>({
    options,
    refreshAccessToken,
    url
  });
}

export const useGetAvailableIntegrationsQuery = () => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();

  return useQuery({
    enabled: !!accessToken && !!namespaceId,
    queryFn: ({ signal }) =>
      fetchInstallableIntegrations({
        accessToken,
        namespaceId,
        signal,
        refreshAccessToken
      }),
    queryKey: ['installableIntegrations']
  });
};
