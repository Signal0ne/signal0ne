import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import Spinner from '../Spinner/Spinner';
import WorkflowSteps from '../WorkflowSteps/WorkflowSteps';
import WorkflowStepDetails from '../WorkflowStepDetails/WorkflowStepDetails';
import './WorkflowsMainPanel.scss';

const WorkflowsMainPanel = () => {
  const { activeWorkflow, isWorkflowLoading } = useWorkflowsContext();

  const getContent = () => {
    if (isWorkflowLoading) return <Spinner />;

    if (!activeWorkflow)
      return (
        <p className="workflows-main-panel--empty">
          Please select the workflow from the side panel.
        </p>
      );

    return (
      <>
        <span className="workflows-breadcrumbs">
          Workflows/{activeWorkflow?.name.replace(/ /g, '-')}
        </span>
        <section className="workflows-workflow">
          <WorkflowStepDetails />
          <WorkflowSteps />
        </section>
      </>
    );
  };

  return <main className="workflows-main-panel">{getContent()}</main>;
};

export default WorkflowsMainPanel;
