import type { RefreshAccessTokenFn } from '../../contexts/AuthProvider/AuthProvider';
import type { Workflow } from '../../data/dummyWorkflows';
import { fetchDataWithAuth } from '../utils';
import { toast } from 'react-toastify';
import { useAuthContext } from '../useAuthContext';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';

type NewWorkflow = Record<string, unknown>;

interface UploadWorkflowProps {
  accessToken: string;
  namespaceId: string;
  newWorkflow: NewWorkflow;
  refreshAccessToken: RefreshAccessTokenFn;
}

interface UploadWorkflowResponse {
  webhook: string;
  workflow: Workflow;
}

interface UseUploadWorkflowMutationProps {
  openModal: () => void;
  setWebhookUrl: (webhookUrl: string) => void;
}

const uploadWorkflow = ({
  accessToken,
  newWorkflow,
  namespaceId,
  refreshAccessToken
}: UploadWorkflowProps) => {
  if (!accessToken || !namespaceId) throw new Error('Something went wrong!');

  const url = `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/workflow/create`;
  const options = {
    body: JSON.stringify(newWorkflow),
    headers: {
      Authorization: `Bearer ${accessToken}}`,
      'Content-Type': 'application/json'
    },
    method: 'POST'
  };

  return fetchDataWithAuth<UploadWorkflowResponse>({
    options,
    refreshAccessToken,
    url
  });
};

export const useUploadWorkflowMutation = ({
  openModal,
  setWebhookUrl
}: UseUploadWorkflowMutationProps) => {
  const { accessToken, namespaceId, refreshAccessToken } = useAuthContext();

  const queryClient = useQueryClient();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: (newWorkflow: NewWorkflow) =>
      uploadWorkflow({
        accessToken,
        newWorkflow,
        namespaceId,
        refreshAccessToken
      }),
    onError: () => {
      toast.error("Couldn't upload selected workflow");
    },
    onSuccess: data => {
      setWebhookUrl(data.webhook);
      navigate(`/${data.workflow.id}`);
      openModal();

      queryClient.invalidateQueries({ queryKey: ['workflows'] });

      toast.success(`Workflow uploaded successfully`);
    }
  });
};
