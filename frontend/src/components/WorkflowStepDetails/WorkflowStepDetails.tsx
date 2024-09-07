import { getIntegrationIcon } from '../../utils/utils';
import { IWorkflowStep, IWorkflowTrigger } from '../../data/dummyWorkflows';
import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import './WorkflowStepDetails.scss';

const WorkflowStepDetails = () => {
  const { activeStep } = useWorkflowsContext();

  const isWorkflowStep = (
    step: IWorkflowStep | IWorkflowTrigger
  ): step is IWorkflowStep => {
    return (step as IWorkflowStep).integration !== undefined;
  };

  if (activeStep && !isWorkflowStep(activeStep)) {
    const isWebhook = 'webhook' in activeStep;
    const triggerOutput = isWebhook
      ? activeStep.webhook.output
      : activeStep.scheduled.output;

    return (
      <div className="workflow-step-details">
        <div className="workflow-step-details-group title">
          <h3 className="workflow-step-details-group-header">
            Type
            <div className="workflow-step-details-group-header-icon">
              {getIntegrationIcon(isWebhook ? 'webhook' : 'scheduled')}
            </div>
          </h3>
          <input
            className="workflow-step-input"
            readOnly
            type="text"
            value={isWebhook ? 'Webhook' : 'Scheduled'}
          />
        </div>
        {!isWebhook && (
          <div className="workflow-step-details-group scheduled">
            <h3 className="workflow-step-details-group-header">Interval</h3>
            <input
              className="workflow-step-input"
              readOnly
              type="text"
              value={activeStep.scheduled.interval}
            />
          </div>
        )}
        <div className="workflow-step-details-group output">
          <h3 className="workflow-step-details-group-header">Output</h3>
          <div className="workflow-step-details-group-content">
            {triggerOutput &&
              Object.entries(triggerOutput).map((step, index) => (
                <div className="workflow-step-group-item" key={step[0] + index}>
                  <span className="workflow-step-group-item-key">
                    {`${step[0]}:{{`}
                  </span>
                  <span className="workflow-step-group-item-value">
                    {step[1]}
                  </span>
                  {'}}'}
                </div>
              ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="workflow-step-details">
      <div className="workflow-step-details-group title">
        <h3 className="workflow-step-details-group-header">
          Title
          <div className="workflow-step-details-group-header-icon">
            {getIntegrationIcon(activeStep?.integrationType || '')}
          </div>
        </h3>
        <input
          className="workflow-step-input"
          readOnly
          type="text"
          value={activeStep?.displayName || ''}
        />
      </div>
      <div className="workflow-step-details-group title">
        <h3 className="workflow-step-details-group-header">Name</h3>
        <input
          className="workflow-step-input"
          readOnly
          type="text"
          value={activeStep?.name || ''}
        />
      </div>
      <div className="workflow-step-details-group function">
        <h3 className="workflow-step-details-group-header">Function</h3>
        <input
          className="workflow-step-input"
          readOnly
          type="text"
          value={activeStep?.function || ''}
        />
      </div>
      <div className="workflow-step-details-group input">
        <h3 className="workflow-step-details-group-header">Input</h3>
        <div className="workflow-step-details-group-content">
          {activeStep?.input &&
            Object.entries(activeStep.input).map((step, index) => (
              <div className="workflow-step-group-item" key={step[0] + index}>
                <span className="workflow-step-group-item-key">
                  {`${step[0]}:{{`}
                </span>
                <span className="workflow-step-group-item-value">
                  {step[1]}
                </span>
                {'}}'}
              </div>
            ))}
        </div>
      </div>
      <div className="workflow-step-details-group output">
        <h3 className="workflow-step-details-group-header">Output</h3>
        <div className="workflow-step-details-group-content">
          {activeStep?.output &&
            Object.entries(activeStep.output).map((step, index) => (
              <div className="workflow-step-group-item" key={step[0] + index}>
                <span className="workflow-step-group-item-key">
                  {`${step[0]}:{{`}
                </span>
                <span className="workflow-step-group-item-value">
                  {step[1]}
                </span>
                {'}}'}
              </div>
            ))}
        </div>
      </div>
      <div className="workflow-step-details-group condition">
        <h3 className="workflow-step-details-group-header">Condition</h3>
        <div className="workflow-step-details-group-content">
          {activeStep?.condition && (
            <div className="workflow-step-group-item">
              {activeStep.condition}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default WorkflowStepDetails;
