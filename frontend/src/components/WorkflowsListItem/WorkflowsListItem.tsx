import { handleKeyDown } from '../../utils/utils';
import { Workflow } from '../../data/dummyWorkflows';
import classNames from 'classnames';
import './WorkflowsListItem.scss';

interface WorkflowsListItemProps {
  isActive?: boolean;
  onClick: (w: Workflow) => void;
  workflow: Workflow;
}

const WorkflowsListItem = ({
  isActive,
  onClick,
  workflow
}: WorkflowsListItemProps) => (
  <li
    className={classNames('workflows-list-item', { active: isActive })}
    key={workflow.id}
    onClick={() => onClick(workflow)}
    onKeyDown={handleKeyDown(() => onClick(workflow))}
    tabIndex={0}
  >
    <h3 className="workflows-list-item--title">{workflow.name}</h3>
    <p className="workflows-list-item--description">{workflow.description}</p>
  </li>
);

export default WorkflowsListItem;
