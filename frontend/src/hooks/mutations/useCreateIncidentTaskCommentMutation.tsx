import type { Incident } from '../../contexts/IncidentsProvider/IncidentsProvider';
import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import { fetchDataWithAuth } from '../utils';
import { toast } from 'react-toastify';
import { useAuthContext } from '../useAuthContext';
import { useIncidentsContext } from '../useIncidentsContext';
import { useMutation } from '@tanstack/react-query';

interface CreateIncidentTaskCommentMutationProps {
  incidentId: string | undefined;
  incidentTaskId: string | undefined;
  newComment: TaskCommentPayload;
}

interface CreateIncidentTaskCommentProps {
  accessToken: string;
  incidentId: string | undefined;
  incidentTaskId: string | undefined;
  namespaceId: string;
  newComment: TaskCommentPayload;
  refreshAccessToken: RefreshAccessTokenFn;
}

interface TaskCommentPayload {
  content: string;
  title: string;
}

interface TaskUpdateResponse {
  updatedIncident: Incident;
}

interface UseCreateIncidentTaskCommentMutationProps {
  handleCloseCommentEditor: () => void;
}

const createTaskComment = async ({
  accessToken,
  incidentId,
  incidentTaskId,
  namespaceId,
  newComment,
  refreshAccessToken
}: CreateIncidentTaskCommentProps) => {
  if (
    !accessToken ||
    !namespaceId ||
    !incidentId ||
    !incidentTaskId ||
    !newComment
  )
    throw new Error('Something went wrong!');

  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/${incidentId}/${incidentTaskId}/add-task-comment`;
  const options = {
    body: JSON.stringify(newComment),
    headers: {
      Authorization: `Bearer ${accessToken}`,
      'Content-Type': 'application/json'
    },
    method: 'POST'
  };

  return fetchDataWithAuth<TaskUpdateResponse>({
    options,
    refreshAccessToken,
    url
  });
};

export const useCreateIncidentTaskCommentMutation = ({
  handleCloseCommentEditor
}: UseCreateIncidentTaskCommentMutationProps) => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();
  const { setSelectedIncident } = useIncidentsContext();

  return useMutation({
    mutationFn: ({
      incidentId,
      incidentTaskId,
      newComment
    }: CreateIncidentTaskCommentMutationProps) =>
      createTaskComment({
        accessToken,
        incidentId,
        incidentTaskId,
        namespaceId,
        newComment,
        refreshAccessToken
      }),
    onError: () => {
      toast.error('Failed to save the comment');
    },
    onSuccess: data => {
      setSelectedIncident(data.updatedIncident);
      handleCloseCommentEditor();
      toast.success('Task comment saved successfully');
    }
  });
};
