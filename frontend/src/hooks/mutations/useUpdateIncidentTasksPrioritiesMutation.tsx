import type {
  IIncidentTask,
  Incident
} from '../../contexts/IncidentsProvider/IncidentsProvider';
import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import { fetchDataWithAuth } from '../utils';
import { toast } from 'react-toastify';
import { useAuthContext } from '../useAuthContext';
import { useIncidentsContext } from '../useIncidentsContext';
import { useMutation } from '@tanstack/react-query';

interface TaskUpdateResponse {
  updatedIncident: Incident;
}

interface UpdateIncidentTaskAssigneeProps {
  accessToken: string;
  incidentId: string | undefined;
  namespaceId: string;
  newTasks: UpdateIncidentTasksPrioritiesPayload;
  refreshAccessToken: RefreshAccessTokenFn;
}

interface UpdateIncidentTaskPriorityMutationProps {
  incidentId: string | undefined;
  newTasks: UpdateIncidentTasksPrioritiesPayload;
}

interface UpdateIncidentTasksPrioritiesPayload {
  incidentTasks: IIncidentTask[];
}

const updateTasksPriorities = async ({
  accessToken,
  incidentId,
  namespaceId,
  newTasks,
  refreshAccessToken
}: UpdateIncidentTaskAssigneeProps) => {
  if (!accessToken || !namespaceId || !incidentId || !newTasks)
    throw new Error('Something went wrong!');

  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/${incidentId}/update-tasks-priority`;
  const options = {
    body: JSON.stringify(newTasks),
    headers: {
      Authorization: `Bearer ${accessToken}`,
      'Content-Type': 'application/json'
    },
    method: 'PATCH'
  };

  return fetchDataWithAuth<TaskUpdateResponse>({
    options,
    refreshAccessToken,
    url
  });
};

export const useUpdateIncidentTasksPrioritiesMutation = () => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();
  const { setSelectedIncident } = useIncidentsContext();

  return useMutation({
    mutationFn: ({
      incidentId,
      newTasks
    }: UpdateIncidentTaskPriorityMutationProps) =>
      updateTasksPriorities({
        accessToken,
        incidentId,
        namespaceId,
        newTasks,
        refreshAccessToken
      }),
    onError: () => {
      toast.error('Failed to update the tasks priorities');
    },
    onSuccess: data => {
      setSelectedIncident(data.updatedIncident);
      toast.success('Tasks priorities updated successfully');
    }
  });
};
