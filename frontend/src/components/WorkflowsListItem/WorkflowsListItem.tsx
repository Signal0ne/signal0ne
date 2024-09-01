import { handleKeyDown } from '../../utils/utils';
import { Workflow } from '../../data/dummyWorkflows';
import classNames from 'classnames';
import './WorkflowsListItem.scss';

interface WorkflowsListItemProps {
  isActive?: boolean;
  onClick: (w: Workflow) => void;
  workflow: Workflow;
}

const DESCRIPTION_DISPLAY_LIMIT = 60;

const WorkflowsListItem = ({
  isActive,
  onClick,
  workflow
}: WorkflowsListItemProps) => {
  const descriptionToDisplay =
    workflow.description.length > DESCRIPTION_DISPLAY_LIMIT
      ? `${workflow.description.slice(0, DESCRIPTION_DISPLAY_LIMIT)}...`
      : workflow.description;

  return (
    <li
      className={classNames('workflows-list-item', { active: isActive })}
      key={workflow.id}
      onClick={() => onClick(workflow)}
      onKeyDown={handleKeyDown(() => onClick(workflow))}
      tabIndex={0}
    >
      <h3 className="workflows-list-item--title">{workflow.name}</h3>
      <p
        className="workflows-list-item--description"
        title={workflow.description}
      >
        {descriptionToDisplay}
      </p>
    </li>
  );
};

export default WorkflowsListItem;
