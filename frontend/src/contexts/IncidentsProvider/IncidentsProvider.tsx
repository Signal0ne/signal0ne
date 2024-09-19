import { createContext, ReactNode, useEffect, useState } from 'react';
import { useAuthContext } from '../../hooks/useAuthContext';
import { toast } from 'react-toastify';

export interface Incident {
  assignee: IncidentAssignee;
  history: string[];
  id: string;
  severity: string;
  summary: string;
  tasks: IIncidentTask[];
  timestamp: number;
  title: string;
}

interface IncidentAssignee {
  email: string;
  id: string;
  name: string;
  photoUrl: string;
  type: string;
}

export interface IncidentsContextType {
  incidents: Incident[];
  isLoading: boolean;
  selectedIncident: Incident | null;
  setIncidents: (incidents: Incident[]) => void;
  setSelectedIncident: (incident: Incident | null) => void;
}

interface IncidentsProviderProps {
  children: ReactNode;
}

export interface IIncidentTask {
  assignee: IncidentAssignee;
  isDone: boolean;
  items: IncidentTaskItem[];
  priority: number;
  taskName: string;
}

export interface IncidentTaskItem {
  content: IncidentTaskItemContent[];
  source: string;
}

interface IncidentTaskItemContent {
  key: string;
  value: string;
  valueKey: 'graph' | 'markdown' | 'text';
}

interface IncidentsResponseBody {
  incidents: Incident[];
}

export const IncidentsContext = createContext<IncidentsContextType | null>(
  null
);

export const IncidentsProvider = ({ children }: IncidentsProviderProps) => {
  const [incidents, setIncidents] = useState<Incident[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [selectedIncident, setSelectedIncident] = useState<Incident | null>(
    null
  );

  const { namespaceId } = useAuthContext();

  useEffect(() => {
    if (!namespaceId) return;

    setIsLoading(true);

    const fetchIncidents = async () => {
      try {
        const response = await fetch(
          `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/incidents`
        );

        if (!response.ok) {
          throw new Error('Failed to fetch incidents');
        }

        const data: IncidentsResponseBody = await response.json();
        setIncidents(data.incidents);
      } catch (error) {
        if (error instanceof Error) {
          toast.error(error.message);
        } else {
          toast.error('An unexpected error occurred. Please try again later.');
        }
      } finally {
        setIsLoading(false);
      }
    };

    fetchIncidents();
  }, [namespaceId]);

  const VALUE = {
    incidents,
    isLoading,
    selectedIncident,
    setIncidents,
    setSelectedIncident
  };

  return (
    <IncidentsContext.Provider value={VALUE}>
      {children}
    </IncidentsContext.Provider>
  );
};
