import {
  getFormattedFormLabel,
  getInputType,
  getIntegrationGradientColor
} from '../../utils/utils';
import { Integration } from '../../contexts/IntegrationsProvider/IntegrationsProvider';
import { SubmitHandler, useForm } from 'react-hook-form';
import { toast } from 'react-toastify';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useIntegrationsContext } from '../../hooks/useIntegrationsContext';
import { useEffect, useMemo, useState } from 'react';
import Input from '../Input/Input';
import ReactModal, { Styles } from 'react-modal';
import Spinner from '../Spinner/Spinner';
import './InstallIntegrationModal.scss';

interface Error {
  message: string;
}

interface FormData {
  [key: string]: unknown;
}

interface GetInstalledIntegrationsResponse {
  installedIntegrations: Integration[];
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

const InstallIntegrationModal = () => {
  const [error, setError] = useState<Error | null>(null);

  const {
    formState: { errors },
    handleSubmit,
    register,
    reset
  } = useForm();

  const { namespaceId } = useAuthContext();
  const {
    isModalOpen,
    selectedIntegration,
    setInstalledIntegrations,
    setIsModalOpen,
    setSelectedIntegration
  } = useIntegrationsContext();

  useEffect(() => {
    setError(null);
  }, [selectedIntegration]);

  const submitForm: SubmitHandler<FormData> = async data => {
    const { name, ...rest } = data;

    if (!selectedIntegration) return;

    const newIntegration = {
      config: rest,
      name,
      type: selectedIntegration.type
    };

    try {
      setError(null);

      let res: Response;

      if (selectedIntegration.id) {
        res = await fetch(
          `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/integration/${selectedIntegration.id}`,
          {
            body: JSON.stringify(newIntegration),
            headers: {
              'Content-Type': 'application/json'
            },
            method: 'PATCH'
          }
        );
      } else {
        res = await fetch(
          `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/integration`,
          {
            body: JSON.stringify(newIntegration),
            headers: {
              'Content-Type': 'application/json'
            },
            method: 'POST'
          }
        );
      }

      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.error);
      }

      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/integration/installed`
      );
      const data: GetInstalledIntegrationsResponse = await response.json();

      setInstalledIntegrations(data.installedIntegrations);
      setIsModalOpen(false);

      toast.success(
        `Integration ${selectedIntegration.id ? 'updated' : 'installed'} successfully`
      );
    } catch (err) {
      toast.error('Failed to install integration');

      if (err instanceof Error) {
        setError(err);
      } else {
        setError(new Error('An unknown error occurred'));
      }
    }
  };

  const formattedSelectedIntegration = useMemo(() => {
    if (!selectedIntegration || !selectedIntegration?.config) return null;

    const formFields = [
      {
        key: 'name',
        value: selectedIntegration.id ? selectedIntegration.name : 'string'
      }
    ];

    if (
      selectedIntegration?.config !== null &&
      'url' in selectedIntegration.config
    ) {
      formFields.push({ key: 'url', value: selectedIntegration.config.url });
    }

    Object.entries(selectedIntegration.config).map(entry => {
      const [key, value] = entry;

      if (key === 'url') return;

      formFields.push({ key, value });
    });

    return formFields;
  }, [selectedIntegration]);

  return (
    <ReactModal
      isOpen={isModalOpen}
      onAfterClose={() => {
        setSelectedIntegration(null);
        reset();
      }}
      onRequestClose={() => {
        setIsModalOpen(false);
      }}
      style={CUSTOM_STYLES}
    >
      {selectedIntegration ? (
        <div className="install-integration-container">
          <h3 className="form-title">
            {selectedIntegration.id ? 'Edit ' : 'Install '}
            <span
              className="integration-name"
              style={{
                backgroundImage: getIntegrationGradientColor(
                  selectedIntegration.type
                )
              }}
            >
              {selectedIntegration.name}
            </span>{' '}
            integration
          </h3>
          <form className="form-content" onSubmit={handleSubmit(submitForm)}>
            {formattedSelectedIntegration &&
              formattedSelectedIntegration.map(entry => {
                const { key, value } = entry;

                const errorMessage =
                  typeof errors[key]?.message === 'string'
                    ? { message: errors[key]?.message }
                    : undefined;

                return (
                  <div className="form-field" key={key}>
                    <Input
                      defaultValue={selectedIntegration?.id ? value : undefined}
                      error={errorMessage}
                      id={`field-${key}`}
                      label={getFormattedFormLabel(key)}
                      placeholder={`Enter ${getFormattedFormLabel(key)} here...`}
                      type={getInputType(key)} //TODO: Find a way to determinate it from the BE
                      {...register(key, {
                        pattern:
                          key === 'url'
                            ? {
                                message: 'Invalid URL address',
                                value: /^https?:\/\/[^\s/$?#].[^\s]*(:\d+)?$/
                              }
                            : undefined,
                        required: 'This field is required'
                      })}
                    />
                  </div>
                );
              })}
            {error ? <p className="error-msg">{error.message}</p> : null}
            <button className="submit" type="submit">
              {selectedIntegration.id ? 'Save Changes' : 'Install'}
            </button>
          </form>
        </div>
      ) : (
        <Spinner />
      )}
    </ReactModal>
  );
};

export default InstallIntegrationModal;
