import type { InstalledIntegration } from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import { getIntegrationIcon, handleKeyDown } from '../../utils/utils';
import { useGetIntegrationByIdMutation } from '../../hooks/mutations/useGetIntegrationByIdMutation';
import { useIntegrationsContext } from '../../hooks/useIntegrationsContext';
import classNames from 'classnames';
import './InstalledIntegrationsListItem.scss';

interface InstalledIntegrationsListItemProps {
  integration: InstalledIntegration;
}

const InstalledIntegrationsListItem = ({
  integration
}: InstalledIntegrationsListItemProps) => {
  const { selectedIntegration } = useIntegrationsContext();

  const { mutate } = useGetIntegrationByIdMutation({
    integrationId: integration.id
  });

  const handleInstalledIntegrationClick = () => mutate();

  return (
    <li
      className={classNames('installed-integrations-list-item', {
        active: selectedIntegration?.id === integration.id
      })}
      onClick={handleInstalledIntegrationClick}
      onKeyDown={handleKeyDown(handleInstalledIntegrationClick)}
      tabIndex={0}
    >
      <div className="integration-icon">
        {getIntegrationIcon(integration.type)}
      </div>
      <span className="integration-name">{integration.name}</span>
    </li>
  );
};

export default InstalledIntegrationsListItem;
