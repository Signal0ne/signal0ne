import { handleKeyDown } from '../../utils/utils';
import { Workflow } from '../../data/dummyWorkflows';
import { useNavigate } from 'react-router-dom';
import classNames from 'classnames';
import './WorkflowsListItem.scss';

interface WorkflowsListItemProps {
  isActive?: boolean;
  workflow: Workflow;
}

const DESCRIPTION_DISPLAY_LIMIT = 60;

const WorkflowsListItem = ({ isActive, workflow }: WorkflowsListItemProps) => {
  const navigate = useNavigate();

  const handleWorkflowClick = () => navigate(`/${workflow.id}`);

  const descriptionToDisplay =
    workflow.description.length > DESCRIPTION_DISPLAY_LIMIT
      ? `${workflow.description.slice(0, DESCRIPTION_DISPLAY_LIMIT)}...`
      : workflow.description;

  return (
    <li
      className={classNames('workflows-list-item', { active: isActive })}
      key={workflow.id}
      onClick={handleWorkflowClick}
      onKeyDown={handleKeyDown(handleWorkflowClick)}
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
