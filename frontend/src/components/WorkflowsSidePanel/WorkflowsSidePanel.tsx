import { ChangeEvent, useEffect, useMemo, useState } from 'react';
import { toast } from 'react-toastify';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import { Workflow } from '../../data/dummyWorkflows';
import FileUploadButton from '../FileUploadButton/FileUploadButton';
import SearchInput from '../SearchInput/SearchInput';
import WorkflowsList from '../WorkflowsList/WorkflowsList';
import './WorkflowsSidePanel.scss';

interface WorkflowsResponseBody {
  workflows: Workflow[];
}

const WorkflowsSidePanel = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [search, setSearch] = useState('');

  const { namespaceId } = useAuthContext();
  const { setWorkflows, workflows } = useWorkflowsContext();

  useEffect(() => {
    if (!namespaceId) return;

    const fetchWorkflows = async () => {
      setIsLoading(true);

      try {
        const response = await fetch(
          `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/workflow/workflows`
        );
        const data: WorkflowsResponseBody = await response.json();

        setWorkflows(data.workflows);
      } catch (error) {
        console.error('Error fetching workflows:', error);
        toast.error("Couldn't fetch workflows");
      } finally {
        setIsLoading(false);
      }
    };

    fetchWorkflows();
  }, [namespaceId, setWorkflows]);

  const handleSearch = (e: ChangeEvent) => {
    const target = e.target as HTMLInputElement;
    setSearch(target.value);
  };

  const FILTERED_WORKFLOWS = useMemo(
    () =>
      workflows.filter(workflow =>
        workflow.name.toLowerCase().includes(search.toLowerCase())
      ),
    [search, workflows]
  );

  return (
    <aside className="workflows-side-panel">
      <div className="workflows-side-panel-title">
        <h1>Workflows</h1>
      </div>
      <div className="workflows-side-panel-content">
        <SearchInput
          onChange={handleSearch}
          placeholder="Search for Workflows..."
          value={search}
        />
        <FileUploadButton />
        <WorkflowsList
          isEmpty={workflows.length === 0}
          isLoading={isLoading}
          workflows={FILTERED_WORKFLOWS}
        />
      </div>
    </aside>
  );
};

export default WorkflowsSidePanel;
