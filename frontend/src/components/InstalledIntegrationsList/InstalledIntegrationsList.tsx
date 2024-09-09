import { checkDisplayScrollOffset } from '../../utils/utils';
import { Integration } from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import { useEffect, useRef, useState } from 'react';
import classNames from 'classnames';
import InstalledIntegrationsListItem from '../InstalledIntegrationsListItem/InstalledIntegrationsListItem';
import Spinner from '../Spinner/Spinner';
import './InstalledIntegrationsList.scss';

interface InstalledIntegrationsListProps {
  installedIntegrations: Integration[];
  isLoading: boolean;
}

const InstalledIntegrationsList = ({
  installedIntegrations,
  isLoading
}: InstalledIntegrationsListProps) => {
  const [shouldDisplayScrollOffset, setShouldDisplayScrollOffset] =
    useState(false);

  const installedIntegrationsListRef = useRef<HTMLUListElement>(null);

  useEffect(() => {
    if (!installedIntegrationsListRef.current) return;

    const shouldDisplayOffset = checkDisplayScrollOffset(
      installedIntegrationsListRef.current
    );

    setShouldDisplayScrollOffset(shouldDisplayOffset);
  }, [installedIntegrations]);

  return (
    <ul
      className={classNames('installed-integrations-list', {
        'scroll-offset': shouldDisplayScrollOffset
      })}
      ref={installedIntegrationsListRef}
    >
      {isLoading ? (
        <Spinner />
      ) : installedIntegrations?.length ? (
        installedIntegrations.map(integration => (
          <InstalledIntegrationsListItem
            integration={integration}
            key={integration.id}
          />
        ))
      ) : (
        <p className="installed-integrations-list--empty">
          No installed integrations found
          <span className="helpful-msg">
            Select desired integrations from the list on the right to install
          </span>
        </p>
      )}
    </ul>
  );
};

export default InstalledIntegrationsList;
