import type { Incident } from '../../contexts/IncidentsProvider/IncidentsProvider';
import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import { fetchDataWithAuth } from '../utils';
import { toast } from 'react-toastify';
import { useAuthContext } from '../useAuthContext';
import { useIncidentsContext } from '../useIncidentsContext';
import { useMutation } from '@tanstack/react-query';

interface TaskUpdateResponse {
  updatedIncident: Incident;
}

interface UpdateIncidentTaskStatusPayload {
  updatedTaskStatus: boolean;
}

interface UpdateIncidentTaskStatusProps {
  accessToken: string;
  incidentId: string | undefined;
  incidentTaskId: string | undefined;
  namespaceId: string;
  newStatus: UpdateIncidentTaskStatusPayload;
  refreshAccessToken: RefreshAccessTokenFn;
}

interface UpdateNewTaskStatusMutationProps {
  incidentId: string | undefined;
  incidentTaskId: string | undefined;
  newStatus: UpdateIncidentTaskStatusPayload;
}

const updateTaskAssignee = async ({
  accessToken,
  incidentId,
  incidentTaskId,
  namespaceId,
  newStatus,
  refreshAccessToken
}: UpdateIncidentTaskStatusProps) => {
  if (
    !accessToken ||
    !namespaceId ||
    !incidentId ||
    !incidentTaskId ||
    !newStatus
  )
    throw new Error('Something went wrong!');

  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/${incidentId}/${incidentTaskId}/status`;
  const options = {
    body: JSON.stringify(newStatus),
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

export const useUpdateIncidentTaskStatusMutation = () => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();
  const { setSelectedIncident } = useIncidentsContext();

  return useMutation({
    mutationFn: ({
      incidentId,
      incidentTaskId,
      newStatus
    }: UpdateNewTaskStatusMutationProps) =>
      updateTaskAssignee({
        accessToken,
        incidentId,
        incidentTaskId,
        namespaceId,
        newStatus,
        refreshAccessToken
      }),
    onError: () => {
      toast.error('Failed to update task status');
    },
    onSuccess: data => {
      setSelectedIncident(data.updatedIncident);
      toast.success('Task status updated successfully');
    }
  });
};
