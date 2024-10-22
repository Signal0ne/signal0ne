import { ChangeEvent, useEffect, useMemo, useState } from 'react';
import { toast } from 'react-toastify';
import { useGetInstalledIntegrationsQuery } from '../../hooks/queries/useGetInstalledIntegrationsQuery';
import InstalledIntegrationsList from '../InstalledIntegrationsList/InstalledIntegrationsList';
import SearchInput from '../SearchInput/SearchInput';
import './InstalledIntegrationsPanel.scss';

const InstalledIntegrationsPanel = () => {
  const [search, setSearch] = useState('');

  const { data, isError, isLoading } = useGetInstalledIntegrationsQuery();

  useEffect(() => {
    if (isError) toast.error('Cannot load installed integrations');
  }, [isError]);

  const handleSearch = (e: ChangeEvent) => {
    const target = e.target as HTMLInputElement;
    setSearch(target.value);
  };

  const FILTERED_INSTALLED_INTEGRATIONS = useMemo(
    () =>
      (data?.installedIntegrations ?? []).filter(integration =>
        integration.name?.toLowerCase().includes(search.trim().toLowerCase())
      ),
    [data, search]
  );

  return (
    <aside className="installed-integrations-container">
      <h3 className="installed-integrations-title">Your Integrations:</h3>
      <SearchInput
        onChange={handleSearch}
        placeholder="Search for Integration..."
        value={search}
      />
      <InstalledIntegrationsList
        installedIntegrations={FILTERED_INSTALLED_INTEGRATIONS}
        isError={isError}
        isLoading={isLoading}
      />
    </aside>
  );
};

export default InstalledIntegrationsPanel;
