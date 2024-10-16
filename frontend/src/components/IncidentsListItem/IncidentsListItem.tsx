import type { Incident } from '../../contexts/IncidentsProvider/IncidentsProvider';
import { handleKeyDown } from '../../utils/utils';
import { useIncidentsContext } from '../../hooks/useIncidentsContext';
import { useNavigate } from 'react-router-dom';
import classNames from 'classnames';
import './IncidentsListItem.scss';

interface IncidentsListItemProps {
  incident: Incident;
}

const IncidentsListItem = ({ incident }: IncidentsListItemProps) => {
  const { selectedIncident } = useIncidentsContext();

  const navigate = useNavigate();

  const handleIncidentClick = () => navigate(`/incidents/${incident.id}`);

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
          high: incident.severity === 'high',
          moderate: incident.severity === 'moderate',
          low: incident.severity === 'low'
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
