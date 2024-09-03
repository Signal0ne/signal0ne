import { AvailableIntegration } from '../../data/dummyAvailableIntegrations';
import {
  getFormattedFormLabel,
  getIntegrationGradientColor
} from '../../utils/utils';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useEffect, useState } from 'react';
import { useIntegrationsContext } from '../../hooks/useIntegrationsContext';
import AvailableIntegrationsList from '../AvailableIntegrationsList/AvailableIntegrationsList';
import Input from '../Input/Input';
import ReactModal, { Styles } from 'react-modal';
import './AvailableIntegrationsPanel.scss';
interface FetchInstallableIntegrationsResponse {
  installableIntegrations: AvailableIntegration[];
}

const CUSTOM_STYLES: Styles = {
  content: {
    backgroundColor: '#383838',
    border: 'none',
    borderRadius: '8px',
    height: 'max-content',
    margin: 'auto',
    padding: '2rem',
    width: 'max-content'
  },
  overlay: {
    backgroundColor: 'rgba(0, 0, 0, 0.5)'
  }
};

const AvailableIntegrationsPanel = () => {
  const [availableIntegrations, setAvailableIntegrations] = useState<
    AvailableIntegration[]
  >([]);

  const { namespaceId } = useAuthContext();

  const { selectedIntegration, setSelectedIntegration } =
    useIntegrationsContext();

  useEffect(() => {
    if (!namespaceId) return;

    const fetchAvailableIntegrations = async () => {
      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/integration/installable`
      );
      const data: FetchInstallableIntegrationsResponse = await response.json();

      setAvailableIntegrations(data.installableIntegrations);
    };

    fetchAvailableIntegrations();
  }, [namespaceId]);

  console.log(selectedIntegration);
  return (
    <main className="available-integrations-container">
      <h3 className="available-integrations-title">Available Integrations:</h3>
      <AvailableIntegrationsList
        availableIntegrations={availableIntegrations}
      />

      <ReactModal
        isOpen={Boolean(selectedIntegration)}
        onRequestClose={() => setSelectedIntegration(null)}
        style={CUSTOM_STYLES}
      >
        <div className="install-integration-container">
          <h3 className="form-title">
            Install{' '}
            <span
              className="integration-name"
              style={{
                backgroundImage: getIntegrationGradientColor(
                  (selectedIntegration as AvailableIntegration)?.typeName
                )
              }}
            >
              {(selectedIntegration as AvailableIntegration)?.displayName}
            </span>{' '}
            integration
          </h3>
          <form
            className="form"
            onSubmit={e => {
              e.preventDefault();
              console.log('submit');
            }}
          >
            {(selectedIntegration as AvailableIntegration)?.config &&
              Object.entries(selectedIntegration?.config).map(entry => {
                const [key, value] = entry;
                return (
                  <div className="form-field" key={key}>
                    <Input
                      // pattern="^https?://(www.)?+(:[0-9]{4,5})?\/.*$"
                      id={`field-${key}`}
                      label={getFormattedFormLabel(key)}
                      onChange={e => console.log(e.target.value)}
                    />
                  </div>
                );
              })}
            <button className="submit" type="submit">
              Install
            </button>
          </form>
        </div>
      </ReactModal>
    </main>
  );
};

export default AvailableIntegrationsPanel;
