import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import { Workflow } from '../../data/dummyWorkflows';
import WorkflowsListItem from '../WorkflowsListItem/WorkflowsListItem';
import './WorkflowsList.scss';

interface WorkflowsListProps {
  workflows: Workflow[];
}

const WorkflowsList = ({ workflows }: WorkflowsListProps) => {
  const { activeWorkflow, setActiveWorkflow } = useWorkflowsContext();
  const onClickHandler = (workflow: Workflow) => {
    console.log('Clicked Workflow: ', workflow);
    setActiveWorkflow(workflow);
  };

  return (
    <ul className="workflows-list">
      {workflows?.length ? (
        workflows.map(workflow => (
          <WorkflowsListItem
            isActive={workflow.id === activeWorkflow?.id}
            key={workflow.id}
            onClick={onClickHandler}
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
