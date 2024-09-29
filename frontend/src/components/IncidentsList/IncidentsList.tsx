import { checkDisplayScrollOffset } from '../../utils/utils';
import { Incident } from '../../contexts/IncidentsProvider/IncidentsProvider';
import { useEffect, useRef, useState } from 'react';
import { useIncidentsContext } from '../../hooks/useIncidentsContext';
import classNames from 'classnames';
import IncidentsListItem from '../IncidentsListItem/IncidentsListItem';
import Spinner from '../Spinner/Spinner';
import './IncidentsList.scss';

interface IncidentsListProps {
  incidentsList: Incident[];
}

const IncidentsList = ({ incidentsList }: IncidentsListProps) => {
  const [shouldDisplayScrollOffset, setShouldDisplayScrollOffset] =
    useState(false);

  const { isIncidentListLoading } = useIncidentsContext();

  const incidentsListRef = useRef<HTMLUListElement>(null);

  useEffect(() => {
    if (!incidentsListRef.current) return;

    const shouldDisplayOffset = checkDisplayScrollOffset(
      incidentsListRef.current
    );

    setShouldDisplayScrollOffset(shouldDisplayOffset);
  }, [incidentsList]);

  return (
    <ul
      className={classNames('incidents-list', {
        'scroll-offset': shouldDisplayScrollOffset
      })}
      ref={incidentsListRef}
    >
      {isIncidentListLoading ? (
        <Spinner />
      ) : incidentsList.length ? (
        incidentsList.map(incident => (
          <IncidentsListItem incident={incident} key={incident.id} />
        ))
      ) : (
        <p className="incidents-list--empty">
          No incidents found
          <br />
          <span className="helpful-msg">Please refine your search</span>
        </p>
      )}
    </ul>
  );
};

export default IncidentsList;
