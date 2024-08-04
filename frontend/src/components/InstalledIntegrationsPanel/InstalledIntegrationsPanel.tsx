import InstalledIntegrationsList from '../InstalledIntegrationsList/InstalledIntegrationsList';
import './InstalledIntegrationsPanel.scss';

const InstalledIntegrationsPanel = () => (
  <aside className="installed-integrations-container">
    <h3 className="installed-integrations-title">Your Integrations:</h3>
    <InstalledIntegrationsList />
  </aside>
);

export default InstalledIntegrationsPanel;
