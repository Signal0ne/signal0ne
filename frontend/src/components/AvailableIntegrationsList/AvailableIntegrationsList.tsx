import { Integration } from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import AvailableIntegrationsListItem from '../AvailableIntegrationsListItem/AvailableIntegrationsListItem';
import './AvailableIntegrationsList.scss';

interface AvailableIntegrationsListProps {
  availableIntegrations: Integration[];
}

const AvailableIntegrationsList = ({
  availableIntegrations
}: AvailableIntegrationsListProps) => (
  <ul className="available-integrations-list">
    {availableIntegrations.map(integration => (
      <AvailableIntegrationsListItem
        integration={integration}
        key={integration.type}
      />
    ))}
  </ul>
);

export default AvailableIntegrationsList;
