import { getIntegrationIcon } from '../../utils/utils';
import { Step } from '../../data/dummyWorkflows';
import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import './WorkflowStep.scss';

interface WorkflowStepProps {
  step: Step;
}

const WorkflowStep = ({ step }: WorkflowStepProps) => {
  const { activeStep, setActiveStep } = useWorkflowsContext();

  const handleStepClick = () => setActiveStep(step);

  return (
    <div
      className={`workflow-step ${
        step.name === activeStep?.name ? 'active' : ''
      }`}
      key={step.name + step.function}
      onClick={handleStepClick}
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
  );
};

export default WorkflowStep;
