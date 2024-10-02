import { getIntegrationIcon, handleKeyDown } from '../../utils/utils';
import { IWorkflowTrigger } from '../../data/dummyWorkflows';
import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import classNames from 'classnames';
import './WorkflowStepTrigger.scss';

interface WorkflowStepProps {
  step: IWorkflowTrigger;
}

const WorkflowStepTrigger = ({ step }: WorkflowStepProps) => {
  const { setActiveStep } = useWorkflowsContext() as {
    setActiveStep: (step: IWorkflowTrigger) => void;
  };

  const handleStepClick = () => setActiveStep(step);

  const isWebhook = 'webhook' in step;
  const stepOptions = isWebhook ? step.webhook.output : step.scheduled.output;

  return (
    <div className="workflow-step-container">
      <p className={classNames('workflow-step-index')}>Trigger</p>
      <div
        className={classNames('workflow-step-content')}
        onClick={handleStepClick}
        onKeyDown={handleKeyDown(handleStepClick)}
        tabIndex={0}
      >
        <div className="workflow-step-icon">
          {getIntegrationIcon(isWebhook ? 'webhook' : 'scheduled')}
        </div>
        <div className="workflow-step-info-container">
          <span className="workflow-step-info-name">
            {isWebhook ? 'Webhook' : 'Scheduled'}
          </span>
          {!isWebhook && (
            <span className="workflow-step-info-function">
              Interval: {step.scheduled.interval}
            </span>
          )}
          {isWebhook && (
            <span className="workflow-step-info-function">
              Integration: {step.webhook.integration}
            </span>
          )}
        </div>
        <div className="workflow-step-output">
          {stepOptions &&
            Object.entries(stepOptions).map((output, index) => {
              if (index === 3)
                return (
                  <span className="more-dots" key={index}>
                    ...
                  </span>
                );
              if (index > 3) return null;

              return (
                <span
                  className="output-content"
                  key={index}
                >{`${output[0]}:{{${output[1]}}}`}</span>
              );
            })}
        </div>
      </div>
    </div>
  );
};

export default WorkflowStepTrigger;
