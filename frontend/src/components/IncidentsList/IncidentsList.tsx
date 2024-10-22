import type { Incident } from '../../contexts/IncidentsProvider/IncidentsProvider';
import { checkDisplayScrollOffset } from '../../utils/utils';
import { useEffect, useRef, useState } from 'react';
import classNames from 'classnames';
import IncidentsListItem from '../IncidentsListItem/IncidentsListItem';
import Spinner from '../Spinner/Spinner';
import './IncidentsList.scss';

interface IncidentsListProps {
  incidentsList: Incident[];
  isError: boolean;
  isLoading: boolean;
}

const IncidentsList = ({
  incidentsList,
  isError,
  isLoading
}: IncidentsListProps) => {
  const [shouldDisplayScrollOffset, setShouldDisplayScrollOffset] =
    useState(false);

  const incidentsListRef = useRef<HTMLUListElement>(null);

  useEffect(() => {
    if (!incidentsListRef.current) return;

    const shouldDisplayOffset = checkDisplayScrollOffset(
      incidentsListRef.current
    );

    setShouldDisplayScrollOffset(shouldDisplayOffset);
  }, [incidentsList]);

  const getContent = () => {
    if (isLoading) return <Spinner />;

    if (isError)
      return (
        <p className="incidents-list--empty">
          Something went wrong!
          <span className="helpful-msg">Please try again later</span>
        </p>
      );

    if (!incidentsList?.length)
      return (
        <p className="incidents-list--empty">
          No incidents found
          <span className="helpful-msg">Please refine your search</span>
        </p>
      );

    return incidentsList.map(incident => (
      <IncidentsListItem incident={incident} key={incident.id} />
    ));
  };

  return (
    <ul
      className={classNames('incidents-list', {
        'scroll-offset': shouldDisplayScrollOffset
      })}
      ref={incidentsListRef}
    >
      {getContent()}
    </ul>
  );
};

export default IncidentsList;
