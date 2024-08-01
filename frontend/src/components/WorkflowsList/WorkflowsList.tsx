import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import { Workflow } from '../../data/dummyWorkflows';
import FileUploadButton from '../FileUploadButton/FileUploadButton';
import WorkflowsListItem from '../WorkflowsListItem/WorkflowsListItem';
import './WorkflowsList.scss';

interface WorkflowsListProps {
  workflows: Workflow[];
}

const WorkflowsList = ({ workflows }: WorkflowsListProps) => {
  const { activeWorkflow, setActiveWorkflow } = useWorkflowsContext();

  const handleListItemClick = (workflow: Workflow) => {
    console.log('Clicked Workflow: ', workflow);
    setActiveWorkflow(workflow);
  };

  return (
    <>
      <FileUploadButton />
      <ul className="workflows-list">
        {workflows?.length ? (
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
    </>
  );
};

export default WorkflowsList;
