import { AvailableIntegration } from '../../data/dummyAvailableIntegrations';
import { getIntegrationIcon, handleKeyDown } from '../../utils/utils';
import { useIntegrationsContext } from '../../hooks/useIntegrationsContext';
import './AvailableIntegrationsListItem.scss';

interface AvailableIntegrationsListItemProps {
  integration: AvailableIntegration;
}

const AvailableIntegrationsListItem = ({
  integration
}: AvailableIntegrationsListItemProps) => {
  const { setSelectedIntegration } = useIntegrationsContext();

  const handleAvailableIntegrationClick = () =>
    setSelectedIntegration(integration);

  return (
    <li
      className="available-integrations-list-item"
      key={integration.name}
      onClick={handleAvailableIntegrationClick}
      onKeyDown={handleKeyDown(handleAvailableIntegrationClick)}
      tabIndex={0}
    >
      <div className="available-integration-icon">
        {getIntegrationIcon(integration.icon)}
      </div>
      <span className="available-integration-name">{integration.name}</span>
    </li>
  );
};

export default AvailableIntegrationsListItem;
