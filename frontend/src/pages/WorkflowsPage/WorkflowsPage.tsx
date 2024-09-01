import { WorkflowsProvider } from '../../contexts/WorkflowsProvider/WorkflowsProvider';
import WorkflowsMainPanel from '../../components/WorkflowsMainPanel/WorkflowsMainPanel';
import WorkflowsSidePanel from '../../components/WorkflowsSidePanel/WorkflowsSidePanel';
import './WorkflowsPage.scss';

const WorkflowsPage = () => (
  <WorkflowsProvider>
    <div className="workflows-container">
      <WorkflowsSidePanel />
      <WorkflowsMainPanel />
    </div>
  </WorkflowsProvider>
);

export default WorkflowsPage;
