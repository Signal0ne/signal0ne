import { AvailableIntegration } from '../../data/dummyAvailableIntegrations';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useEffect, useState } from 'react';
import AvailableIntegrationsList from '../AvailableIntegrationsList/AvailableIntegrationsList';
import './AvailableIntegrationsPanel.scss';

interface FetchInstallableIntegrationsResponse {
  installableIntegrations: AvailableIntegration[];
}

const AvailableIntegrationsPanel = () => {
  const [availableIntegrations, setAvailableIntegrations] = useState<
    AvailableIntegration[]
  >([]);

  const { namespaceId } = useAuthContext();

  useEffect(() => {
    if (!namespaceId) return;

    const fetchAvailableIntegrations = async () => {
      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/integration/installable`
      );
      const data: FetchInstallableIntegrationsResponse = await response.json();

      setAvailableIntegrations(data.installableIntegrations);
    };

    fetchAvailableIntegrations();
  }, [namespaceId]);

  return (
    <main className="available-integrations-container">
      <h3 className="available-integrations-title">Available Integrations:</h3>
      <AvailableIntegrationsList
        availableIntegrations={availableIntegrations}
      />
    </main>
  );
};

export default AvailableIntegrationsPanel;
