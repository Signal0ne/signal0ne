import {
  AvailableIntegration,
  DUMMY_AVAILABLE_INTEGRATIONS
} from '../../data/dummyAvailableIntegrations';
import { createContext, ReactNode, useState } from 'react';
import {
  DUMMY_INSTALLED_INTEGRATIONS,
  InstalledIntegration
} from '../../data/dummyInstalledIntegrations';

export interface IntegrationsContextType {
  availableIntegrations: AvailableIntegration[];
  installedIntegrations: InstalledIntegration[];
  selectedIntegration: AvailableIntegration | InstalledIntegration | null;
  setSelectedIntegration: (
    integration: AvailableIntegration | InstalledIntegration
  ) => void;
}

interface IntegrationsProviderProps {
  children: ReactNode;
}

const IntegrationsContext = createContext<IntegrationsContextType | undefined>(
  undefined
);

export const IntegrationsProvider = ({
  children
}: IntegrationsProviderProps) => {
  const [selectedIntegration, setSelectedIntegration] = useState<
    AvailableIntegration | InstalledIntegration | null
  >(null);

  const VALUE = {
    availableIntegrations: DUMMY_AVAILABLE_INTEGRATIONS,
    installedIntegrations: DUMMY_INSTALLED_INTEGRATIONS,
    selectedIntegration,
    setSelectedIntegration
  };

  return (
    <IntegrationsContext.Provider value={VALUE}>
      {children}
    </IntegrationsContext.Provider>
  );
};

export default IntegrationsContext;
