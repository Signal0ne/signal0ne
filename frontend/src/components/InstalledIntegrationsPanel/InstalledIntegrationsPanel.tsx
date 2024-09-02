import { InstalledIntegration } from '../../data/dummyInstalledIntegrations';
import { PlusIcon } from '../Icons/Icons';
import { toast } from 'react-toastify';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useEffect, useState } from 'react';
import { useIntegrationsContext } from '../../hooks/useIntegrationsContext';
import InstalledIntegrationsList from '../InstalledIntegrationsList/InstalledIntegrationsList';
import './InstalledIntegrationsPanel.scss';

interface GetInstalledIntegrationsResponse {
  installedIntegrations: InstalledIntegration[];
}

const InstalledIntegrationsPanel = () => {
  const [isLoading, setIsLoading] = useState(true);

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
        console.error('Error fetching installed integrations:', error);
        toast.error('Cannot fetch installed integrations');
      } finally {
        setIsLoading(false);
      }
    };

    fetchInstalledIntegrations();
  }, [namespaceId, setInstalledIntegrations]);

  return (
    <aside className="installed-integrations-container">
      <h3 className="installed-integrations-title">Your Integrations:</h3>
      <button className="install-integration-btn" type="button">
        <PlusIcon height={24} width={24} />
        Install Integration
      </button>
      <InstalledIntegrationsList
        isLoading={isLoading}
        installedIntegrations={installedIntegrations}
      />
    </aside>
  );
};

export default InstalledIntegrationsPanel;
