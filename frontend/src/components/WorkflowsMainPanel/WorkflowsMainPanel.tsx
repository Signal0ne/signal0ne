import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import './WorkflowsMainPanel.scss';

const WorkflowsMainPanel = () => {
  const { activeWorkflow } = useWorkflowsContext();

  return (
    <main className="workflows-main-panel">
      {activeWorkflow ? (
        <>
          <span className="workflows-breadcrumbs">
            Workflows/{activeWorkflow.name.replace(/ /g, '-')}
          </span>
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
