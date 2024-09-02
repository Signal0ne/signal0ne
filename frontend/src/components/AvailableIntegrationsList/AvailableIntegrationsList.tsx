import { AvailableIntegration } from '../../data/dummyAvailableIntegrations';
import AvailableIntegrationsListItem from '../AvailableIntegrationsListItem/AvailableIntegrationsListItem';
import './AvailableIntegrationsList.scss';

interface AvailableIntegrationsListProps {
  availableIntegrations: AvailableIntegration[];
}

const AvailableIntegrationsList = ({
  availableIntegrations
}: AvailableIntegrationsListProps) => {
  // const { availableIntegrations } = useIntegrationsContext();

  return (
    <ul className="available-integrations-list">
      {availableIntegrations.map(integration => (
        <AvailableIntegrationsListItem
          integration={integration}
          key={integration.typeName}
        />
      ))}
    </ul>
  );
};

export default AvailableIntegrationsList;
