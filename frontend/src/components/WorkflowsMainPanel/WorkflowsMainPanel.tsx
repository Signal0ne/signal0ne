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
          <section className="workflows-workflow">
            <div className="workflow-details">
              <div className="workflow-details-group title">
                <h3 className="workflow-details-group-header">Title</h3>
                <input
                  className="workflow-input"
                  type="text"
                  value={activeWorkflow.name}
                />
              </div>
              <div className="workflow-details-group function-name">
                <h3 className="workflow-details-group-header">Function</h3>
                <input
                  className="workflow-input"
                  readOnly
                  type="text"
                  value={activeWorkflow.function}
                />
              </div>
              <div className="workflow-details-group input">
                <h3 className="workflow-details-group-header">Input</h3>
                <div className="workflow-input-container"></div>
              </div>
              <div className="workflow-details-group output">
                <h3 className="workflow-details-group-header">Output</h3>
                <div className="workflow-output-container"></div>
              </div>
              <div className="workflow-details-group condition">
                <h3 className="workflow-details-group-header">Condition</h3>
                <div className="workflow-condition-container"></div>
              </div>
            </div>
            <div className="workflow-steps-container"></div>
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
