import { Integration } from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useEffect, useState } from 'react';
import AvailableIntegrationsList from '../AvailableIntegrationsList/AvailableIntegrationsList';
import InstallIntegrationModal from '../InstallIntegrationModal/InstallIntegrationModal';
import './AvailableIntegrationsPanel.scss';

interface FetchInstallableIntegrationsResponse {
  installableIntegrations: Integration[];
}

const AvailableIntegrationsPanel = () => {
  const [availableIntegrations, setAvailableIntegrations] = useState<
    Integration[]
  >([]);

  const { accessToken, namespaceId } = useAuthContext();

  useEffect(() => {
    if (!namespaceId || !accessToken) return;

    const fetchAvailableIntegrations = async () => {
      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/integration/installable`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`
          }
        }
      );
      const data: FetchInstallableIntegrationsResponse = await response.json();

      setAvailableIntegrations(data.installableIntegrations);
    };

    fetchAvailableIntegrations();
  }, [accessToken, namespaceId]);

  return (
    <main className="available-integrations-container">
      <h3 className="available-integrations-title">Available Integrations:</h3>
      <AvailableIntegrationsList
        availableIntegrations={availableIntegrations}
      />
      <InstallIntegrationModal />
    </main>
  );
};

export default AvailableIntegrationsPanel;
