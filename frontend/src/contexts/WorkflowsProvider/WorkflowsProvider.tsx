import { createContext, ReactNode, useState } from 'react';
import {
  IWorkflowStep,
  IWorkflowTrigger,
  Workflow
} from '../../data/dummyWorkflows';

export interface WorkflowsContextType {
  activeStep: IWorkflowStep | IWorkflowTrigger | null;
  activeWorkflow: Workflow | null;
  setActiveStep: (step: IWorkflowStep | IWorkflowTrigger | null) => void;
  setActiveWorkflow: (workflow: any) => void;
  setWorkflows: (workflows: Workflow[]) => void;
  workflows: Workflow[];
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
  const [workflows, setWorkflows] = useState<Workflow[]>([]);

  const VALUE = {
    activeStep,
    activeWorkflow,
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
