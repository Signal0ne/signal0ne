import { WorkflowsProvider } from '../../contexts/WorkflowsProvider/WorkflowsProvider';
import WorkflowsMainPanel from '../../components/WorkflowsMainPanel/WorkflowsMainPanel';
import WorkflowsSidePanel from '../../components/WorkflowsSidePanel/WorkflowsSidePanel';
import './WorkflowsPage.scss';

const WorkflowsPage = () => (
  <div className="workflows-container">
    <WorkflowsProvider>
      <WorkflowsSidePanel />
      <WorkflowsMainPanel />
    </WorkflowsProvider>
  </div>
);

export default WorkflowsPage;
