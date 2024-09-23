import { useIncidentsContext } from '../../hooks/useIncidentsContext';
import { useMemo, useState } from 'react';
import IncidentsList from '../IncidentsList/IncidentsList';
import SearchInput from '../SearchInput/SearchInput';
import './IncidentsSidebar.scss';

const IncidentsSidebar = () => {
  const [searchValue, setSearchValue] = useState('');

  const { incidents } = useIncidentsContext();

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchValue(e.target.value);
  };

  const FILTERED_INCIDENTS = useMemo(
    () =>
      incidents?.filter(incident =>
        incident.title?.toLowerCase().includes(searchValue.trim().toLowerCase())
      ),
    [incidents, searchValue]
  );

  return (
    <aside className="incidents-sidebar">
      <h3 className="incidents-title">Incidents</h3>
      <SearchInput
        onChange={handleSearch}
        placeholder="Search for Incident..."
        value={searchValue}
      />
      <IncidentsList incidentsList={FILTERED_INCIDENTS} />
    </aside>
  );
};

export default IncidentsSidebar;
