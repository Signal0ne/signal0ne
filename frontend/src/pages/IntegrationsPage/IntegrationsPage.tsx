import { IntegrationsProvider } from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import AvailableIntegrationsPanel from '../../components/AvailableIntegrationsPanel/AvailableIntegrationsPanel';
import InstalledIntegrationsPanel from '../../components/InstalledIntegrationsPanel/InstalledIntegrationsPanel';
import './IntegrationsPage.scss';

const IntegrationsPage = () => (
  <div className="integrations-page">
    <IntegrationsProvider>
      <InstalledIntegrationsPanel />
      <AvailableIntegrationsPanel />
    </IntegrationsProvider>
  </div>
);

export default IntegrationsPage;
