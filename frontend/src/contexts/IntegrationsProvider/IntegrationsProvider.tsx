import { createContext, ReactNode, useState } from 'react';

export interface InstalledIntegration extends Integration {
  id: string;
}

export interface Integration {
  config: Record<string, string> | null;
  id?: string;
  imageUri: string;
  name: string;
  type: string;
}

export interface IntegrationsContextType {
  isModalOpen: boolean;
  selectedIntegration: Integration | null;
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
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedIntegration, setSelectedIntegration] =
    useState<Integration | null>(null);

  const VALUE = {
    isModalOpen,
    selectedIntegration,
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
