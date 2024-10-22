import type {
  IWorkflowStep,
  IWorkflowTrigger,
  Workflow
} from '../../data/dummyWorkflows';
import { createContext, ReactNode, useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { useGetWorkflowByIdQuery } from '../../hooks/queries/useGetWorkflowByIdQuery';
import { useParams } from 'react-router-dom';

export interface WorkflowsContextType {
  activeStep: IWorkflowStep | IWorkflowTrigger | null;
  activeWorkflow: Workflow | null;
  isWorkflowLoading: boolean;
  setActiveStep: (step: IWorkflowStep | IWorkflowTrigger | null) => void;
}

interface WorkflowsProviderProps {
  children: ReactNode;
}

export const WorkflowsContext = createContext<WorkflowsContextType | undefined>(
  undefined
);

export const WorkflowsProvider = ({ children }: WorkflowsProviderProps) => {
  const [activeStep, setActiveStep] =
    useState<WorkflowsContextType['activeStep']>(null);
  const [activeWorkflow, setActiveWorkflow] = useState<Workflow | null>(null);

  const { workflowId } = useParams<{ workflowId?: string }>();

  const { data, isError, isLoading } = useGetWorkflowByIdQuery(workflowId);

  useEffect(() => {
    if (isError) toast.error('Failed to get workflow data');
  }, [isError]);

  useEffect(() => {
    if (data?.workflow) setActiveWorkflow(data.workflow);
  }, [data]);

  useEffect(() => {
    if (!workflowId) setActiveWorkflow(null);
  }, [workflowId]);

  useEffect(() => {
    activeWorkflow && setActiveStep(activeWorkflow?.steps[0]);
  }, [activeWorkflow]);

  const VALUE = {
    activeStep,
    activeWorkflow,
    isWorkflowLoading: isLoading,
    setActiveStep
  };

  return (
    <WorkflowsContext.Provider value={VALUE}>
      {children}
    </WorkflowsContext.Provider>
  );
};
