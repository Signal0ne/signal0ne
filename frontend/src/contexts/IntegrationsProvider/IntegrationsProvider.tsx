
import { createContext, ReactNode, useState } from 'react';

export interface Integration {
  config: Record<string, string> | null;
  id?: string;
  imageUri: string;
  name: string;
  type: string;
}

export interface IntegrationsContextType {
  availableIntegrations: Integration[];
  installedIntegrations: Integration[];
  selectedIntegration: Integration | null;
  setAvailableIntegrations: (integrations: Integration[]) => void;
  setInstalledIntegrations: (integrations: Integration[]) => void;
  setSelectedIntegration: (integration: Integration | null) => void;
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
    Integration[]
  >([]);
  const [installedIntegrations, setInstalledIntegrations] = useState<
    Integration[]
  >([]);
  const [selectedIntegration, setSelectedIntegration] =
    useState<Integration | null>(null);

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
