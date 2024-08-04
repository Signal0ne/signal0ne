import {
  ArrowDown,
  BackStageIcon,
  JaegerIcon,
  PrometheusIcon,
  SlackIcon
} from '../Icons/Icons';
import { ReactNode, useEffect } from 'react';
import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import './WorkflowsMainPanel.scss';

const getIntegrationIcon = (integrationName: string) => {
  const icons: Record<string, ReactNode> = {
    backstage: <BackStageIcon />,
    jaeger: <JaegerIcon />,
    prometheus: <PrometheusIcon />,
    slack: <SlackIcon />
  };

  return icons[integrationName] || null;
};

const WorkflowsMainPanel = () => {
  const { activeWorkflow, activeStep, setActiveStep } = useWorkflowsContext();

  useEffect(() => {
    if (!activeWorkflow?.steps[1]) return;
    setActiveStep(activeWorkflow?.steps[1]);
  }, [activeWorkflow, setActiveStep]);

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

  console.log(activeWorkflow, activeStep);
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
                <h3 className="workflow-details-group-header">
                  Title
                  <div className="workflow-details-group-header-icon">
                    {getIntegrationIcon(activeStep?.integration || '')}
                  </div>
                </h3>
                <input
                  className="workflow-input"
                  readOnly
                  type="text"
                  value={activeStep?.name || ''}
                />
              </div>
              <div className="workflow-details-group function">
                <h3 className="workflow-details-group-header">Function</h3>
                <input
                  className="workflow-input"
                  readOnly
                  type="text"
                  value={activeStep?.function || ''}
                />
              </div>
              <div className="workflow-details-group input">
                <h3 className="workflow-details-group-header">Input</h3>
                <div className="workflow-details-group-content input">
                  {activeStep?.input &&
                    Object.entries(activeStep.input).map((step, index) => (
                      <div key={index} className="workflow-group-item">
                        <span className="workflow-group-item-key">
                          {`${step[0]}:{{`}
                        </span>
                        <span className="workflow-group-item-value">
                          {step[1]}
                        </span>
                        {'}}'}
                      </div>
                    ))}
                </div>
              </div>
              <div className="workflow-details-group output">
                <h3 className="workflow-details-group-header">Output</h3>
                <div className="workflow-details-group-content output">
                  {activeStep?.output &&
                    Object.entries(activeStep.output).map((step, index) => (
                      <div key={index} className="workflow-group-item">
                        <span className="workflow-group-item-key">
                          {`${step[0]}:{{`}
                        </span>
                        <span className="workflow-group-item-value">
                          {step[1]}
                        </span>
                        {'}}'}
                      </div>
                    ))}
                </div>
              </div>
              <div className="workflow-details-group condition">
                <h3 className="workflow-details-group-header">Condition</h3>
                <div className="workflow-details-group-content condition">
                  {activeStep?.condition && (
                    <div className="workflow-group-item">
                      {activeStep.condition}
                    </div>
                  )}
                </div>
              </div>
            </div>
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
                    <div
                      className={`workflow-step ${
                        step.name === activeStep?.name ? 'active' : ''
                      }`}
                      key={index}
                      onClick={() => setActiveStep(step)}
                    >
                      <div className="workflow-step-icon">
                        {getIntegrationIcon(step.integration)}
                      </div>
                      <div className="workflow-step-info-container">
                        <span className="workflow-step-info-name">
                          {step.name}
                        </span>
                        <span className="workflow-step-info-function">
                          Function: {step.function}
                        </span>
                      </div>
                      <div className="workflow-step-output">
                        {step?.output &&
                          Object.entries(step.output).map((output, index) => {
                            return (
                              <span
                                key={index}
                              >{`${output[0]}:{{${output[1]}}}`}</span>
                            );
                          })}
                      </div>
                    </div>
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
