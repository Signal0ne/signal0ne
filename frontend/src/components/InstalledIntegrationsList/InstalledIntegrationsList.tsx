import { useIntegrationsContext } from '../../hooks/useIntegrationsContext';
import InstalledIntegrationsListItem from '../InstalledIntegrationsListItem/InstalledIntegrationsListItem';
import './InstalledIntegrationsList.scss';

const InstalledIntegrationsList = () => {
  const { installedIntegrations } = useIntegrationsContext();

  return (
    <ul className="installed-integrations-list">
      {installedIntegrations.map(integration => (
        <InstalledIntegrationsListItem
          integration={integration}
          key={integration.name}
        />
      ))}
      {/* TODO: remove after connecting the Backend */}
      {installedIntegrations.map(integration => (
        <InstalledIntegrationsListItem
          integration={integration}
          key={integration.name}
        />
      ))}
      {installedIntegrations.map(integration => (
        <InstalledIntegrationsListItem
          integration={integration}
          key={integration.name}
        />
      ))}
    </ul>
  );
};

export default InstalledIntegrationsList;
