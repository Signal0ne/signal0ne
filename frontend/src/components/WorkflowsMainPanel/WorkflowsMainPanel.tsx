import { useEffect } from 'react';
import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import './WorkflowsMainPanel.scss';
import WorkflowStepDetails from '../WorkflowStepDetails/WorkflowStepDetails';
import WorkflowSteps from '../WorkflowSteps/WorkflowSteps';

const WorkflowsMainPanel = () => {
  const { activeWorkflow, activeStep, setActiveStep } = useWorkflowsContext();

  useEffect(() => {
    if (!activeWorkflow?.steps[1]) return;
    setActiveStep(activeWorkflow?.steps[1]);
  }, [activeWorkflow, setActiveStep]);

  console.log(activeWorkflow, activeStep);

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
