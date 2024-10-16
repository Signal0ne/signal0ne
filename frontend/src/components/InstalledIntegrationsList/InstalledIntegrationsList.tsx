import type { InstalledIntegration } from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import { checkDisplayScrollOffset } from '../../utils/utils';
import { useEffect, useRef, useState } from 'react';
import classNames from 'classnames';
import InstalledIntegrationsListItem from '../InstalledIntegrationsListItem/InstalledIntegrationsListItem';
import Spinner from '../Spinner/Spinner';
import './InstalledIntegrationsList.scss';

interface InstalledIntegrationsListProps {
  installedIntegrations: InstalledIntegration[];
  isError: boolean;
  isLoading: boolean;
}

const InstalledIntegrationsList = ({
  installedIntegrations,
  isError,
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

  const getContent = () => {
    if (isLoading) return <Spinner />;

    if (isError)
      return (
        <p className="installed-integrations-list--empty">
          Something went wrong!
          <span className="helpful-msg">Please try again later.</span>
        </p>
      );

    if (!installedIntegrations?.length)
      return (
        <p className="installed-integrations-list--empty">
          No integrations found
          <span className="helpful-msg">Please refine your search</span>
        </p>
      );

    return installedIntegrations.map(integration => (
      <InstalledIntegrationsListItem
        integration={integration}
        key={integration.id}
      />
    ));
  };

  return (
    <ul
      className={classNames('installed-integrations-list', {
        'scroll-offset': shouldDisplayScrollOffset
      })}
      ref={installedIntegrationsListRef}
    >
      {getContent()}
    </ul>
  );
};

export default InstalledIntegrationsList;
