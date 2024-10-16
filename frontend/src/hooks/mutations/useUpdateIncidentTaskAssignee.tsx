import type {
  Incident,
  IncidentAssignee
} from '../../contexts/IncidentsProvider/IncidentsProvider';
import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import { fetchDataWithAuth } from '../utils';
import { toast } from 'react-toastify';
import { useAuthContext } from '../useAuthContext';
import { useIncidentsContext } from '../useIncidentsContext';
import { useMutation } from '@tanstack/react-query';

interface TaskAssigneePayload {
  assignee: IncidentAssignee | undefined;
}

interface TaskUpdateResponse {
  updatedIncident: Incident;
}

interface UpdateIncidentTaskAssigneeProps {
  accessToken: string;
  incidentId: string | undefined;
  incidentTaskId: string | undefined;
  namespaceId: string;
  newAssignee: TaskAssigneePayload;
  refreshAccessToken: RefreshAccessTokenFn;
}

interface UpdateNewTaskAssigneeMutationProps {
  incidentId: string | undefined;
  incidentTaskId: string | undefined;
  newAssignee: TaskAssigneePayload;
}

const updateTaskAssignee = async ({
  accessToken,
  incidentId,
  incidentTaskId,
  namespaceId,
  newAssignee,
  refreshAccessToken
}: UpdateIncidentTaskAssigneeProps) => {
  if (
    !accessToken ||
    !namespaceId ||
    !incidentId ||
    !incidentTaskId ||
    !newAssignee
  )
    throw new Error('Something went wrong!');

  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/${incidentId}/${incidentTaskId}/assignee`;
  const options = {
    body: JSON.stringify(newAssignee),
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

export const useUpdateIncidentTaskAssigneeMutation = () => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();
  const { setSelectedIncident } = useIncidentsContext();

  return useMutation({
    mutationFn: ({
      incidentId,
      incidentTaskId,
      newAssignee
    }: UpdateNewTaskAssigneeMutationProps) =>
      updateTaskAssignee({
        accessToken,
        incidentId,
        incidentTaskId,
        namespaceId,
        newAssignee,
        refreshAccessToken
      }),
    onError: () => {
      toast.error('Failed to update task assignee');
    },
    onSuccess: data => {
      setSelectedIncident(data.updatedIncident);
      toast.success('Task assignee updated successfully');
    }
  });
};
