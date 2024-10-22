import { toast } from 'react-toastify';
import { useEffect } from 'react';
import { useGetAvailableIntegrationsQuery } from '../../hooks/queries/useGetAvailableIntegrationsQuery';
import AvailableIntegrationsList from '../AvailableIntegrationsList/AvailableIntegrationsList';
import InstallIntegrationModal from '../InstallIntegrationModal/InstallIntegrationModal';
import Spinner from '../Spinner/Spinner';
import './AvailableIntegrationsPanel.scss';

const AvailableIntegrationsPanel = () => {
  const { data, isError, isLoading } = useGetAvailableIntegrationsQuery();

  useEffect(() => {
    if (isError) toast.error('Cannot load available integrations');
  }, [isError]);

  const getContent = () => {
    if (isLoading) return <Spinner />;

    if (isError)
      return (
        <p className="error-msg">
          Something went wrong!
          <span className="error-msg-subtext">Please try again later.</span>
        </p>
      );

    return (
      <>
        <AvailableIntegrationsList
          availableIntegrations={data?.installableIntegrations ?? []}
        />
        <InstallIntegrationModal />
      </>
    );
  };

  return (
    <main className="available-integrations-container">
      <h3 className="available-integrations-title">Available Integrations:</h3>
      {getContent()}
    </main>
  );
};

export default AvailableIntegrationsPanel;
