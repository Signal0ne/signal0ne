import { ChangeEvent, useEffect, useMemo, useState } from 'react';
import { DUMMY_WORKFLOWS } from '../../data/dummyWorkflows';
import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import SearchInput from '../SearchInput/SearchInput';
import WorkflowsList from '../WorkflowsList/WorkflowsList';
import './WorkflowsSidePanel.scss';

const WorkflowsSidePanel = () => {
  const [search, setSearch] = useState('');

  const { setWorkflows, workflows } = useWorkflowsContext();

  useEffect(() => {
    //TODO: Fetch workflows from API
    setTimeout(() => setWorkflows(DUMMY_WORKFLOWS), 1000);
  }, [setWorkflows]);

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
        <WorkflowsList workflows={FILTERED_WORKFLOWS} />
      </div>
    </aside>
  );
};

export default WorkflowsSidePanel;
