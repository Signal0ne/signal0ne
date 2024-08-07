import { useIntegrationsContext } from '../../hooks/useIntegrationsContext';
import AvailableIntegrationsListItem from '../AvailableIntegrationsListItem/AvailableIntegrationsListItem';
import './AvailableIntegrationsList.scss';

const AvailableIntegrationsList = () => {
  const { availableIntegrations } = useIntegrationsContext();

  return (
    <ul className="available-integrations-list">
      {availableIntegrations.map(integration => (
        <AvailableIntegrationsListItem
          integration={integration}
          key={integration.name}
        />
      ))}
      {/* TODO: remove after connecting the Backend */}
      {availableIntegrations.map(integration => (
        <AvailableIntegrationsListItem
          integration={integration}
          key={integration.name}
        />
      ))}
      {availableIntegrations.map(integration => (
        <AvailableIntegrationsListItem
          integration={integration}
          key={integration.name}
        />
      ))}
    </ul>
  );
};

export default AvailableIntegrationsList;
