import { useContext } from 'react';
import { WorkflowsContext } from '../contexts/WorkflowsProvider/WorkflowsProvider';

export const useWorkflowsContext = () => {
  const context = useContext(WorkflowsContext);

  if (!context) {
    throw new Error(
      'useWorkflowsContext must be used within a WorkflowsProvider'
    );
  }

  return context;
};
