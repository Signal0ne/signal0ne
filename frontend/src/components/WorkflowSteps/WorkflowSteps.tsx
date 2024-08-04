import { ArrowDown } from '../Icons/Icons';
import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import WorkflowStep from '../WorkflowStep/WorkflowStep';
import './WorkflowSteps.scss';

const WorkflowSteps = () => {
  const { activeWorkflow } = useWorkflowsContext();

  const calcStepsListHeight = () => {
    const workflowsContainer =
      document.querySelector('.workflows-workflow')?.getBoundingClientRect()
        .height ?? 0;
    const workflowInfoContainer =
      document
        .querySelector('.workflow-info-container')
        ?.getBoundingClientRect().height ?? 0;

    if (!workflowsContainer) return '100%';

    return workflowsContainer - workflowInfoContainer;
  };

  if (!activeWorkflow) return null;

  return (
    <div className="workflow-steps-container">
      <div className="workflow-info-container">
        <h3 className="workflow-info-name">{activeWorkflow.name}</h3>
        <h5 className="workflow-info-description">
          {activeWorkflow.description}
        </h5>
      </div>
      <div
        className="workflow-steps-list"
        style={{ height: calcStepsListHeight() }}
      >
        {activeWorkflow.steps.map((step, index) => (
          <>
            <WorkflowStep step={step} />
            {index !== activeWorkflow.steps.length - 1 && (
              <ArrowDown
                className="workflow-step-separator"
                height={36}
                width={36}
              />
            )}
          </>
        ))}
      </div>
    </div>
  );
};

export default WorkflowSteps;
