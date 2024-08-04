import { useEffect } from 'react';
import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import './WorkflowsMainPanel.scss';
import WorkflowSteps from '../WorkflowSteps/WorkflowSteps';
import WorkflowStepDetails from '../WorkflowStepDetails/WorkflowStepDetails';

const WorkflowsMainPanel = () => {
  const { activeWorkflow, setActiveStep } = useWorkflowsContext();

  useEffect(() => {
    if (!activeWorkflow?.steps[1]) return;
    setActiveStep(activeWorkflow?.steps[1]);
  }, [activeWorkflow, setActiveStep]);

  return (
    <main className="workflows-main-panel">
      {activeWorkflow ? (
        <>
          <span className="workflows-breadcrumbs">
            Workflows/{activeWorkflow?.name.replace(/ /g, '-')}
          </span>
          <section className="workflows-workflow">
            <WorkflowStepDetails />
            <WorkflowSteps />
          </section>
        </>
      ) : (
        <p className="workflows-main-panel--empty">
          Please select the workflow from the side panel.
        </p>
      )}
    </main>
  );
};

export default WorkflowsMainPanel;
