import { ChangeEvent, useEffect, useMemo, useState } from 'react';
import { Integration } from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import { toast } from 'react-toastify';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useIntegrationsContext } from '../../hooks/useIntegrationsContext';
import InstalledIntegrationsList from '../InstalledIntegrationsList/InstalledIntegrationsList';
import SearchInput from '../SearchInput/SearchInput';
import './InstalledIntegrationsPanel.scss';

interface GetInstalledIntegrationsResponse {
  installedIntegrations: Integration[];
}

const InstalledIntegrationsPanel = () => {
  const [isLoading, setIsLoading] = useState(true);
  const [search, setSearch] = useState('');

  const { namespaceId } = useAuthContext();
  const { installedIntegrations, setInstalledIntegrations } =
    useIntegrationsContext();

  useEffect(() => {
    if (!namespaceId) return;

    const fetchInstalledIntegrations = async () => {
      setIsLoading(true);

      try {
        const response = await fetch(
          `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/integration/installed`
        );

        const data: GetInstalledIntegrationsResponse = await response.json();

        setInstalledIntegrations(data.installedIntegrations);
      } catch (error) {
        toast.error('Cannot fetch installed integrations');
      } finally {
        setIsLoading(false);
      }
    };

    fetchInstalledIntegrations();
  }, [namespaceId, setInstalledIntegrations]);

  const handleSearch = (e: ChangeEvent) => {
    const target = e.target as HTMLInputElement;
    setSearch(target.value);
  };

  const FILTERED_INSTALLED_INTEGRATIONS = useMemo(
    () =>
      installedIntegrations.filter(integration =>
        integration.name?.toLowerCase().includes(search.trim().toLowerCase())
      ),
    [installedIntegrations, search]
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
        isLoading={isLoading}
        installedIntegrations={FILTERED_INSTALLED_INTEGRATIONS}
      />
    </aside>
  );
};

export default InstalledIntegrationsPanel;
