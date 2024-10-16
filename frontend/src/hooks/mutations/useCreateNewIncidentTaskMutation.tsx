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

interface CreateNewIncidentTaskMutationProps {
  incidentId: string;
  newIncidentTask: NewIncidentTaskPayload;
}

interface CreateNewIncidentTaskProps {
  accessToken: string;
  currentIncidentId: string;
  namespaceId: string;
  newIncidentTask: NewIncidentTaskPayload;
  refreshAccessToken: RefreshAccessTokenFn;
}

interface NewIncidentTaskPayload extends Omit<IIncidentTask, 'id'> {}

interface NewIncidentTaskResponse {
  updatedIncident: Incident;
}

interface UseCreateNewIncidentTaskMutationProps {
  handleTaskModalClose: () => void;
}

const createNewIncidentTask = async ({
  accessToken,
  currentIncidentId,
  namespaceId,
  newIncidentTask,
  refreshAccessToken
}: CreateNewIncidentTaskProps) => {
  if (!accessToken || !namespaceId) throw new Error('Something went wrong!');

  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/${currentIncidentId}/tasks`;
  const options = {
    body: JSON.stringify(newIncidentTask),
    headers: {
      Authorization: `Bearer ${accessToken}`,
      'Content-Type': 'application/json'
    },
    method: 'POST'
  };

  return fetchDataWithAuth<NewIncidentTaskResponse>({
    options,
    refreshAccessToken,
    url
  });
};

export const useCreateNewIncidentTaskMutation = ({
  handleTaskModalClose
}: UseCreateNewIncidentTaskMutationProps) => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();
  const { setSelectedIncident } = useIncidentsContext();

  return useMutation({
    mutationFn: ({
      incidentId,
      newIncidentTask
    }: CreateNewIncidentTaskMutationProps) =>
      createNewIncidentTask({
        accessToken,
        currentIncidentId: incidentId,
        namespaceId,
        newIncidentTask,
        refreshAccessToken
      }),
    onError: () => {
      toast.error('Failed to add task');
    },
    onSuccess: data => {
      setSelectedIncident(data.updatedIncident);
      handleTaskModalClose();
      toast.success('Task added successfully');
    }
  });
};
