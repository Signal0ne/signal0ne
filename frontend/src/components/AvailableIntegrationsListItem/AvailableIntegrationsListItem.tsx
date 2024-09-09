import { getIntegrationIcon, handleKeyDown } from '../../utils/utils';
import { Integration } from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import { useIntegrationsContext } from '../../hooks/useIntegrationsContext';
import './AvailableIntegrationsListItem.scss';

interface AvailableIntegrationsListItemProps {
  integration: Integration;
}

const AvailableIntegrationsListItem = ({
  integration
}: AvailableIntegrationsListItemProps) => {
  const { setIsModalOpen, setSelectedIntegration } = useIntegrationsContext();

  const handleAvailableIntegrationClick = () => {
    setIsModalOpen(true);
    setSelectedIntegration(integration);
  };

  return (
    <li
      className="available-integrations-list-item"
      key={integration.type}
      onClick={handleAvailableIntegrationClick}
      onKeyDown={handleKeyDown(handleAvailableIntegrationClick)}
      tabIndex={0}
    >
      <div className="available-integration-icon">
        {getIntegrationIcon(integration.type)}
      </div>
      <span className="available-integration-name">{integration.name}</span>
    </li>
  );
};

export default AvailableIntegrationsListItem;
