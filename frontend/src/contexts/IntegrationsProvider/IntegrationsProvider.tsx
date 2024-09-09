import { createContext, ReactNode, useState } from 'react';

export interface Integration {
  config: Record<string, string> | null;
  id?: string;
  imageUri: string;
  name: string;
  type: string;
}

export interface IntegrationsContextType {
  installedIntegrations: Integration[];
  isModalOpen: boolean;
  selectedIntegration: Integration | null;
  setInstalledIntegrations: (integrations: Integration[]) => void;
  setIsModalOpen: (isOpen: boolean) => void;
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
  const [installedIntegrations, setInstalledIntegrations] = useState<
    Integration[]
  >([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedIntegration, setSelectedIntegration] =
    useState<Integration | null>(null);

  const VALUE = {
    installedIntegrations,
    isModalOpen,
    selectedIntegration,
    setInstalledIntegrations,
    setIsModalOpen,
    setSelectedIntegration
  };

  return (
    <IntegrationsContext.Provider value={VALUE}>
      {children}
    </IntegrationsContext.Provider>
  );
};

export default IntegrationsContext;
