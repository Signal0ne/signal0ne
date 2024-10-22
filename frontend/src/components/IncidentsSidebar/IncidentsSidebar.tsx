import { toast } from 'react-toastify';
import { useEffect, useMemo, useState } from 'react';
import { useGetIncidentsQuery } from '../../hooks/queries/useGetIncidentsQuery';
import IncidentsList from '../IncidentsList/IncidentsList';
import SearchInput from '../SearchInput/SearchInput';
import './IncidentsSidebar.scss';

const IncidentsSidebar = () => {
  const [searchValue, setSearchValue] = useState('');

  const { data, isError, isLoading } = useGetIncidentsQuery();

  useEffect(() => {
    if (isError) toast.error("Couldn't get your incidents");
  }, [isError]);

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchValue(e.target.value);
  };

  const FILTERED_INCIDENTS = useMemo(
    () =>
      (data?.incidents ?? [])?.filter(incident =>
        incident.title?.toLowerCase().includes(searchValue.trim().toLowerCase())
      ),
    [data, searchValue]
  );

  return (
    <aside className="incidents-sidebar">
      <h3 className="incidents-title">Incidents</h3>
      <SearchInput
        onChange={handleSearch}
        placeholder="Search for Incident..."
        value={searchValue}
      />
      <IncidentsList
        incidentsList={FILTERED_INCIDENTS}
        isError={isError}
        isLoading={isLoading}
      />
    </aside>
  );
};

export default IncidentsSidebar;
