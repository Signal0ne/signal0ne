import { checkDisplayScrollOffset } from '../../utils/utils';
import { InstalledIntegration } from '../../data/dummyInstalledIntegrations';
import { useEffect, useRef, useState } from 'react';
import classNames from 'classnames';
import InstalledIntegrationsListItem from '../InstalledIntegrationsListItem/InstalledIntegrationsListItem';
import Spinner from '../Spinner/Spinner';
import './InstalledIntegrationsList.scss';

interface InstalledIntegrationsListProps {
  installedIntegrations: InstalledIntegration[];
  isLoading: boolean;
}

const InstalledIntegrationsList = ({
  isLoading,
  installedIntegrations
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
            key={integration.name}
          />
        ))
      ) : (
        <p className="installed-integrations-list--empty">
          No installed integrations found
          <span className="helpful-msg">
            Click the button above to install one
          </span>
        </p>
      )}
    </ul>
  );
};

export default InstalledIntegrationsList;
