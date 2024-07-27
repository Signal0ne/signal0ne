import { useContext } from 'react';
import {
  WorkflowsContext,
  WorkflowsContextType
} from '../contexts/WorkflowsContext/WorkflowsContext';

export const useWorkflowsContext = () => {
  const context = useContext<WorkflowsContextType | undefined>(
    WorkflowsContext
  );

  if (!context) {
    throw new Error(
      'useWorkflowsContext must be used within a WorkflowsProvider'
    );
  }

  return context;
};
