import { createContext, useState } from 'react';
import { Workflow } from '../../data/dummyWorkflows';

export interface WorkflowsContextType {
  activeWorkflow: Workflow | null;
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
  const [activeWorkflow, setActiveWorkflow] = useState<Workflow | null>(null);
  const [workflows, setWorkflows] = useState<Workflow[]>([]);

  const VALUE = {
    activeWorkflow,
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
