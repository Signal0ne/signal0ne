import { handleKeyDown } from '../../utils/utils';
import { Incident } from '../../contexts/IncidentsProvider/IncidentsProvider';
import { toast } from 'react-toastify';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useIncidentsContext } from '../../hooks/useIncidentsContext';
import classNames from 'classnames';
import './IncidentsListItem.scss';

interface IncidentResponse {
  incident: Incident;
}

interface IncidentsListItemProps {
  incident: Incident;
}

const IncidentsListItem = ({ incident }: IncidentsListItemProps) => {
  const { namespaceId } = useAuthContext();
  const { selectedIncident, setIsIncidentPreviewLoading, setSelectedIncident } =
    useIncidentsContext();

  const handleIncidentClick = async () => {
    setIsIncidentPreviewLoading(true);

    try {
      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/${incident.id}`
      );

      if (!response.ok) throw new Error('Failed to fetch incident details');

      const data: IncidentResponse = await response.json();
      setSelectedIncident(data.incident);
    } catch (error) {
      if (error instanceof Error) {
        toast.error(error.message);
      } else {
        toast.error('An unexpected error occurred. Please try again later.');
      }
    } finally {
      setIsIncidentPreviewLoading(false);
    }
  };

  return (
    <li
      className={classNames('incidents-list-item', {
        active: selectedIncident?.id === incident.id
      })}
      onClick={handleIncidentClick}
      onKeyDown={handleKeyDown(handleIncidentClick)}
      tabIndex={0}
    >
      <span
        className={classNames('incidents-list-item-severity', {
          critical: incident.severity === 'critical',
          error: incident.severity === 'error',
          warning: incident.severity === 'warning'
        })}
      />
      <div className="incidents-list-item-info">
        <span className="incidents-list-item-title">{incident.title}</span>
        <span className="incidents-list-item-time">
          {new Date(incident.timestamp * 1000).toLocaleDateString('en-US', {
            day: '2-digit',
            month: 'short',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
          })}
        </span>
      </div>
    </li>
  );
};

export default IncidentsListItem;
