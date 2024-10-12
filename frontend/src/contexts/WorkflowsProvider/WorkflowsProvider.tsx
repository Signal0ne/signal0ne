import { createContext, ReactNode, useEffect, useState } from 'react';
import {
  IWorkflowStep,
  IWorkflowTrigger,
  Workflow
} from '../../data/dummyWorkflows';
import { toast } from 'react-toastify';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useParams } from 'react-router-dom';

export interface WorkflowsContextType {
  activeStep: IWorkflowStep | IWorkflowTrigger | null;
  activeWorkflow: Workflow | null;
  isWorkflowLoading: boolean;
  setActiveStep: (step: IWorkflowStep | IWorkflowTrigger | null) => void;
  setActiveWorkflow: (workflow: any) => void;
  setWorkflows: (workflows: Workflow[]) => void;
  workflows: Workflow[];
}

interface WorkflowsProviderProps {
  children: ReactNode;
}

interface WorkflowResponseBody {
  workflow: Workflow;
}

export const WorkflowsContext = createContext<WorkflowsContextType | undefined>(
  undefined
);

export const WorkflowsProvider = ({ children }: WorkflowsProviderProps) => {
  const [activeStep, setActiveStep] =
    useState<WorkflowsContextType['activeStep']>(null);
  const [activeWorkflow, setActiveWorkflow] = useState<Workflow | null>(null);
  const [isWorkflowLoading, setIsWorkflowLoading] = useState(false);
  const [workflows, setWorkflows] = useState<Workflow[]>([]);

  const { accessToken, namespaceId } = useAuthContext();
  const { workflowId } = useParams<{ workflowId?: string }>();

  useEffect(() => {
    if (!namespaceId || !accessToken) return;

    const fetchWorkflow = async () => {
      try {
        setIsWorkflowLoading(true);

        const response = await fetch(
          `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/workflow/${workflowId}`,
          {
            headers: {
              Authorization: `Bearer ${accessToken}`
            }
          }
        );

        if (!response.ok) throw new Error('Failed to fetch workflow');

        const data: WorkflowResponseBody = await response.json();

        setActiveWorkflow(data.workflow);
      } catch (error) {
        if (error instanceof Error) {
          toast.error(error.message);
        } else {
          toast.error('An unexpected error occurred. Please try again later.');
        }
      } finally {
        setIsWorkflowLoading(false);
      }
    };

    if (workflowId) {
      setActiveWorkflow(null);
      fetchWorkflow();
    } else {
      setActiveWorkflow(null);
    }
  }, [accessToken, namespaceId, workflowId]);

  const VALUE = {
    activeStep,
    activeWorkflow,
    isWorkflowLoading,
    setActiveStep,
    setActiveWorkflow,
    setWorkflows,
    workflows
  };

  return (
    <WorkflowsContext.Provider value={VALUE}>
      {children}
    </WorkflowsContext.Provider>
  );
};
