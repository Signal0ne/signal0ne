import { getIntegrationIcon, handleKeyDown } from '../../utils/utils';
import { Integration } from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import { toast } from 'react-toastify';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useIntegrationsContext } from '../../hooks/useIntegrationsContext';
import classNames from 'classnames';
import './InstalledIntegrationsListItem.scss';

interface InstalledIntegrationsListItemProps {
  integration: Integration;
}

interface InstalledIntegrationResponse {
  integration: Integration;
}

const InstalledIntegrationsListItem = ({
  integration
}: InstalledIntegrationsListItemProps) => {
  const { accessToken, namespaceId } = useAuthContext();
  const { selectedIntegration, setIsModalOpen, setSelectedIntegration } =
    useIntegrationsContext();

  const handleInstalledIntegrationClick = async () => {
    if (!namespaceId || !accessToken) return;

    setIsModalOpen(true);

    try {
      const res = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/integration/${integration.id}`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`
          }
        }
      );

      if (!res.ok) {
        throw new Error(
          'Failed to fetch integration data, please try again later'
        );
      }

      const data: InstalledIntegrationResponse = await res.json();
      setSelectedIntegration(data.integration);
    } catch (err) {
      if (err instanceof Error) {
        toast.error(err.message);
      } else {
        toast.error('Oops! Something went wrong, please try again later');
      }
    }
  };

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
