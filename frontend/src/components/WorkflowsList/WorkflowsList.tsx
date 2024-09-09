import { checkDisplayScrollOffset } from '../../utils/utils';
import { useEffect, useRef, useState } from 'react';
import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import { Workflow } from '../../data/dummyWorkflows';
import classNames from 'classnames';
import Spinner from '../Spinner/Spinner';
import WorkflowsListItem from '../WorkflowsListItem/WorkflowsListItem';
import './WorkflowsList.scss';

interface WorkflowsListProps {
  isEmpty: boolean;
  isLoading: boolean;
  workflows: Workflow[];
}

const WorkflowsList = ({
  isEmpty,
  isLoading,
  workflows
}: WorkflowsListProps) => {
  const [shouldDisplayScrollOffset, setShouldDisplayScrollOffset] =
    useState(false);

  const workflowsListRef = useRef<HTMLUListElement>(null);

  const { activeWorkflow, setActiveWorkflow } = useWorkflowsContext();

  useEffect(() => {
    if (!workflowsListRef.current) return;

    const shouldDisplayOffset = checkDisplayScrollOffset(
      workflowsListRef.current
    );

    setShouldDisplayScrollOffset(shouldDisplayOffset);
  }, [workflows]);

  const handleListItemClick = (workflow: Workflow) =>
    setActiveWorkflow(workflow);

  return (
    <ul
      className={classNames('workflows-list', {
        'scroll-offset': shouldDisplayScrollOffset
      })}
      ref={workflowsListRef}
    >
      {isLoading ? (
        <Spinner />
      ) : workflows?.length ? (
        workflows.map(workflow => (
          <WorkflowsListItem
            isActive={workflow.id === activeWorkflow?.id}
            key={workflow.id}
            onClick={handleListItemClick}
            workflow={workflow}
          />
        ))
      ) : (
        <p className="workflows-list--empty">
          No workflows found
          <span className="helpful-msg">
            {isEmpty
              ? 'Click the button above to upload one'
              : 'Please refine your search'}
          </span>
        </p>
      )}
    </ul>
  );
};

export default WorkflowsList;
