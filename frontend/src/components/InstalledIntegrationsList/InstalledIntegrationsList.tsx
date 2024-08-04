import { useIntegrationsContext } from '../../hooks/useIntegrationsContext';
import InstalledIntegrationsListItem from '../InstalledIntegrationsListItem/InstalledIntegrationsListItem';
import './InstalledIntegrationsList.scss';

const InstalledIntegrationsList = () => {
  const { installedIntegrations } = useIntegrationsContext();

  return (
    <ul className="installed-integrations-list">
      {installedIntegrations.map(integration => (
        <InstalledIntegrationsListItem
          key={integration.name}
          integration={integration}
        />
      ))}
      {/* TODO: remove after connecting the Backend */}
      {installedIntegrations.map(integration => (
        <InstalledIntegrationsListItem
          key={integration.name}
          integration={integration}
        />
      ))}
      {installedIntegrations.map(integration => (
        <InstalledIntegrationsListItem
          key={integration.name}
          integration={integration}
        />
      ))}
    </ul>
  );
};

export default InstalledIntegrationsList;
