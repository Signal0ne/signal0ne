import { ChangeEvent, useEffect, useMemo, useState } from 'react';
import { toast } from 'react-toastify';
import { useGetWorkflowsQuery } from '../../hooks/queries/useGetWorkflowsQuery';
import FileUploadButton from '../FileUploadButton/FileUploadButton';
import SearchInput from '../SearchInput/SearchInput';
import WorkflowsList from '../WorkflowsList/WorkflowsList';
import './WorkflowsSidePanel.scss';

const WorkflowsSidePanel = () => {
  const [search, setSearch] = useState('');

  const { data, isError, isLoading } = useGetWorkflowsQuery();

  useEffect(() => {
    if (isError) toast.error("Couldn't fetch workflows");
  }, [isError]);

  const handleSearch = (e: ChangeEvent<HTMLInputElement>) =>
    setSearch(e.target.value);

  const FILTERED_WORKFLOWS = useMemo(
    () =>
      (data?.workflows ?? []).filter(workflow =>
        workflow.name?.toLowerCase().includes(search.trim().toLowerCase())
      ),
    [data, search]
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
          isEmpty={data?.workflows.length === 0}
          isError={isError}
          isLoading={isLoading}
          workflows={FILTERED_WORKFLOWS}
        />
      </div>
    </aside>
  );
};

export default WorkflowsSidePanel;
