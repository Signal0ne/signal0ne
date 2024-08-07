import { PlusIcon } from '../Icons/Icons';
import InstalledIntegrationsList from '../InstalledIntegrationsList/InstalledIntegrationsList';
import './InstalledIntegrationsPanel.scss';

const InstalledIntegrationsPanel = () => (
  <aside className="installed-integrations-container">
    <h3 className="installed-integrations-title">Your Integrations:</h3>
    <button className="install-integration-btn" type="button">
      <PlusIcon height={24} width={24} />
      Install Integration
    </button>
    <InstalledIntegrationsList />
  </aside>
);

export default InstalledIntegrationsPanel;
