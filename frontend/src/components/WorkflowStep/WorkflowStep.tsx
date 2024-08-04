import { getIntegrationIcon, handleKeyDown } from '../../utils/utils';
import { Step } from '../../data/dummyWorkflows';
import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import classNames from 'classnames';
import './WorkflowStep.scss';

interface WorkflowStepProps {
  index: number;
  step: Step;
}

const WorkflowStep = ({ index, step }: WorkflowStepProps) => {
  const { activeStep, setActiveStep } = useWorkflowsContext();

  const handleStepClick = () => setActiveStep(step);

  return (
    <div className="workflow-step-container">
      <p
        className={classNames('workflow-step-index', {
          active: step.name === activeStep?.name
        })}
      >
        Step: {index + 1}
      </p>
      <div
        className={classNames('workflow-step-content', {
          active: step.name === activeStep?.name
        })}
        key={step.name + index}
        onClick={handleStepClick}
        onKeyDown={handleKeyDown(handleStepClick)}
        tabIndex={0}
      >
        <div className="workflow-step-icon">
          {getIntegrationIcon(step.integration)}
        </div>
        <div className="workflow-step-info-container">
          <span className="workflow-step-info-name">{step.name}</span>
          <span className="workflow-step-info-function">
            Function: {step.function}
          </span>
        </div>
        <div className="workflow-step-output">
          {step?.output &&
            Object.entries(step.output).map((output, index) => {
              return <span key={index}>{`${output[0]}:{{${output[1]}}}`}</span>;
            })}
        </div>
      </div>
    </div>
  );
};

export default WorkflowStep;
