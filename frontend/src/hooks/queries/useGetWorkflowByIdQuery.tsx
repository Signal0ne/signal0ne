import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import type { Workflow } from '../../data/dummyWorkflows';
import { fetchDataWithAuth } from '../utils';
import { useAuthContext } from '../useAuthContext';
import { useQuery } from '@tanstack/react-query';

interface FetchWorkflowByIdProps {
  accessToken: string;
  namespaceId: string;
  refreshAccessToken: RefreshAccessTokenFn;
  signal: AbortSignal;
  workflowId: string | undefined;
}

interface WorkflowResponseBody {
  workflow: Workflow;
}

const fetchWorkflowById = ({
  accessToken,
  namespaceId,
  refreshAccessToken,
  signal,
  workflowId
}: FetchWorkflowByIdProps) => {
  if (!workflowId) return null;

  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/workflow/${workflowId}`;
  const options = {
    headers: {
      Authorization: `Bearer ${accessToken}`
    },
    signal
  };

  return fetchDataWithAuth<WorkflowResponseBody>({
    options,
    refreshAccessToken,
    url
  });
};

export const useGetWorkflowByIdQuery = (workflowId: string | undefined) => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();

  return useQuery({
    enabled: !!accessToken && !!namespaceId && !!workflowId,
    queryFn: ({ signal }) =>
      fetchWorkflowById({
        accessToken,
        namespaceId,
        refreshAccessToken,
        signal,
        workflowId
      }),
    queryKey: ['workflowId', workflowId]
  });
};
