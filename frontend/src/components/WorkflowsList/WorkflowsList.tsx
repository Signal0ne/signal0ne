import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import { Workflow } from '../../data/dummyWorkflows';
import Spinner from '../Spinner/Spinner';
import WorkflowsListItem from '../WorkflowsListItem/WorkflowsListItem';
import './WorkflowsList.scss';

interface WorkflowsListProps {
  isLoading: boolean;
  workflows: Workflow[];
}

const WorkflowsList = ({ isLoading, workflows }: WorkflowsListProps) => {
  const { activeWorkflow, setActiveWorkflow } = useWorkflowsContext();

  const handleListItemClick = (workflow: Workflow) =>
    setActiveWorkflow(workflow);

  return (
    <ul className="workflows-list">
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
        <p className="workflows-list--empty">No workflows found</p>
      )}
    </ul>
  );
};

export default WorkflowsList;
