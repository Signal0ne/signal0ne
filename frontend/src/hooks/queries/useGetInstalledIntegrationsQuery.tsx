import type { InstalledIntegration } from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import { fetchDataWithAuth } from '../utils';
import { useAuthContext } from '../useAuthContext';
import { useQuery } from '@tanstack/react-query';

interface FetchInstalledIntegrationsProps {
  accessToken: string;
  namespaceId: string;
  refreshAccessToken: RefreshAccessTokenFn;
  signal: AbortSignal;
}

interface FetchInstalledIntegrationsResponse {
  installedIntegrations: InstalledIntegration[];
}

function fetchInstalledIntegrations({
  accessToken,
  namespaceId,
  refreshAccessToken,
  signal
}: FetchInstalledIntegrationsProps) {
  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/integration/installed`;
  const options = {
    headers: {
      Authorization: `Bearer ${accessToken}`
    },
    signal
  };

  return fetchDataWithAuth<FetchInstalledIntegrationsResponse>({
    options,
    refreshAccessToken,
    url
  });
}

export const useGetInstalledIntegrationsQuery = () => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();

  return useQuery({
    enabled: !!accessToken && !!namespaceId,
    queryFn: ({ signal }) =>
      fetchInstalledIntegrations({
        accessToken,
        namespaceId,
        signal,
        refreshAccessToken
      }),
    queryKey: ['installedIntegrations']
  });
};
