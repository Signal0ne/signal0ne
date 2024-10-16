import type { Incident } from '../../contexts/IncidentsProvider/IncidentsProvider';
import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import { fetchDataWithAuth } from '../utils';
import { useAuthContext } from '../useAuthContext';
import { useQuery } from '@tanstack/react-query';

interface FetchIncidentProps {
  accessToken: string;
  incidentId: string | undefined;
  namespaceId: string;
  refreshAccessToken: RefreshAccessTokenFn;
  signal: AbortSignal;
}

interface IncidentResponseBody {
  incident: Incident;
}

const fetchIncidentById = ({
  accessToken,
  incidentId,
  namespaceId,
  refreshAccessToken,
  signal
}: FetchIncidentProps) => {
  if (!incidentId) return null;

  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/${incidentId}`;
  const options = {
    headers: {
      Authorization: `Bearer ${accessToken}`
    },
    signal
  };

  return fetchDataWithAuth<IncidentResponseBody>({
    options,
    refreshAccessToken,
    url
  });
};

export const useGetIncidentByIdQuery = (incidentId: string | undefined) => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();

  return useQuery({
    enabled: !!accessToken && !!namespaceId,
    queryFn: ({ signal }) =>
      fetchIncidentById({
        accessToken,
        incidentId,
        namespaceId,
        refreshAccessToken,
        signal
      }),
    queryKey: ['incidentId', incidentId]
  });
};
