import { createContext, useState } from 'react';
import { Step, Workflow } from '../../data/dummyWorkflows';

export interface WorkflowsContextType {
  activeStep: Step | null;
  activeWorkflow: Workflow | null;
  setActiveStep: (step: Step) => void;
  setActiveWorkflow: (workflow: Workflow) => void;
  setWorkflows: (workflows: Workflow[]) => void;
  workflows: Workflow[];
}

export const WorkflowsContext = createContext<WorkflowsContextType | undefined>(
  undefined
);

export const WorkflowsProvider = ({
  children
}: {
  children: React.ReactNode;
}) => {
  const [activeStep, setActiveStep] = useState<Step | null>(null);
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
