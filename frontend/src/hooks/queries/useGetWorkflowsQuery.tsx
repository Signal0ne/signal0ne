import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import type { Workflow } from '../../data/dummyWorkflows';
import { fetchDataWithAuth } from '../utils';
import { useAuthContext } from '../useAuthContext';
import { useQuery } from '@tanstack/react-query';

interface FetchWorkflowsProps {
  accessToken: string;
  namespaceId: string;
  refreshAccessToken: RefreshAccessTokenFn;
  signal: AbortSignal;
}

interface WorkflowsResponseBody {
  workflows: Workflow[];
}

const fetchWorkflows = ({
  accessToken,
  namespaceId,
  refreshAccessToken,
  signal
}: FetchWorkflowsProps) => {
  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/workflow/workflows`;
  const options = {
    headers: {
      Authorization: `Bearer ${accessToken}`
    },
    signal
  };

  return fetchDataWithAuth<WorkflowsResponseBody>({
    options,
    refreshAccessToken,
    url
  });
};

export const useGetWorkflowsQuery = () => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();

  return useQuery({
    enabled: !!accessToken && !!namespaceId,
    queryFn: ({ signal }) =>
      fetchWorkflows({
        accessToken,
        namespaceId,
        refreshAccessToken,
        signal
      }),
    queryKey: ['workflows']
  });
};
