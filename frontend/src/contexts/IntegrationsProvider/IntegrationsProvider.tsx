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
  setAvailableIntegrations: (integrations: AvailableIntegration[]) => void;
  setInstalledIntegrations: (integrations: InstalledIntegration[]) => void;
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
  const [availableIntegrations, setAvailableIntegrations] = useState<
    AvailableIntegration[]
  >(DUMMY_AVAILABLE_INTEGRATIONS);
  const [installedIntegrations, setInstalledIntegrations] = useState<
    InstalledIntegration[]
  >(DUMMY_INSTALLED_INTEGRATIONS);
  const [selectedIntegration, setSelectedIntegration] = useState<
    AvailableIntegration | InstalledIntegration | null
  >(null);

  const VALUE = {
    availableIntegrations,
    installedIntegrations,
    selectedIntegration,
    setAvailableIntegrations,
    setInstalledIntegrations,
    setSelectedIntegration
  };

  return (
    <IntegrationsContext.Provider value={VALUE}>
      {children}
    </IntegrationsContext.Provider>
  );
};

export default IntegrationsContext;
