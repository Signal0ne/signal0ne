import type { Incident } from '../../contexts/IncidentsProvider/IncidentsProvider';
import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import { fetchDataWithAuth } from '../utils';
import { useAuthContext } from '../useAuthContext';
import { useQuery } from '@tanstack/react-query';

interface FetchIncidentsProps {
  accessToken: string;
  namespaceId: string;
  refreshAccessToken: RefreshAccessTokenFn;
  signal: AbortSignal;
}

interface IncidentsResponseBody {
  incidents: Incident[];
}

const fetchIncidents = ({
  accessToken,
  namespaceId,
  refreshAccessToken,
  signal
}: FetchIncidentsProps) => {
  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/incidents`;
  const options = {
    headers: {
      Authorization: `Bearer ${accessToken}`
    },
    signal
  };

  return fetchDataWithAuth<IncidentsResponseBody>({
    options,
    refreshAccessToken,
    url
  });
};

export const useGetIncidentsQuery = () => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();

  return useQuery({
    enabled: !!accessToken && !!namespaceId,
    queryFn: ({ signal }) =>
      fetchIncidents({
        accessToken,
        namespaceId,
        refreshAccessToken,
        signal
      }),
    queryKey: ['incidents']
  });
};
