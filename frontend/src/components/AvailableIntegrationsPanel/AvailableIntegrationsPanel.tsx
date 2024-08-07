import AvailableIntegrationsList from '../AvailableIntegrationsList/AvailableIntegrationsList';
import './AvailableIntegrationsPanel.scss';

const AvailableIntegrationsPanel = () => (
  <main className="available-integrations-container">
    <h3 className="available-integrations-title">Available Integrations:</h3>
    <AvailableIntegrationsList />
  </main>
);

export default AvailableIntegrationsPanel;
