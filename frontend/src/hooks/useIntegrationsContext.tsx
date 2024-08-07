import { useContext } from 'react';
import IntegrationsContext from '../contexts/IntegrationsProvider/IntegrationsProvider';

export const useIntegrationsContext = () => {
  const context = useContext(IntegrationsContext);

  if (!context) {
    throw new Error(
      'useIntegrationsContext must be used within a IntegrationsProvider'
    );
  }

  return context;
};
