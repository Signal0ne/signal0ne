import type {
  Incident,
  IncidentAssignee
} from '../../contexts/IncidentsProvider/IncidentsProvider';
import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import { fetchDataWithAuth } from '../utils';
import { useAuthContext } from '../useAuthContext';
import { useQuery } from '@tanstack/react-query';

interface FetchUsersProps {
  accessToken: string;
  namespaceId: string;
  refreshAccessToken: RefreshAccessTokenFn;
  signal: AbortSignal;
}

interface UsersResponseBody {
  users: IncidentAssignee[];
}

const fetchUsers = ({
  accessToken,
  namespaceId,
  refreshAccessToken,
  signal
}: FetchUsersProps) => {
  const url = `${import.meta.env.VITE_SERVER_API_URL}/namespace/${namespaceId}/users`;
  const options = {
    headers: {
      Authorization: `Bearer ${accessToken}`
    },
    signal
  };

  return fetchDataWithAuth<UsersResponseBody>({
    options,
    refreshAccessToken,
    url
  });
};

export const useGetNamespaceUsersQuery = (
  selectedIncident: Incident | null
) => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();

  return useQuery({
    enabled: !!accessToken && !!namespaceId && !!selectedIncident,
    queryFn: ({ signal }) =>
      fetchUsers({
        accessToken,
        namespaceId,
        refreshAccessToken,
        signal
      }),
    queryKey: ['users']
  });
};
