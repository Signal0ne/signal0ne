import { createContext, ReactNode, useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useParams } from 'react-router-dom';

type IncidentSeverity = 'critical' | 'high' | 'moderate' | 'low';
export interface Incident {
  assignee: IncidentAssignee;
  history: string[];
  id: string;
  severity: IncidentSeverity;
  summary: string;
  tasks: IIncidentTask[];
  timestamp: number;
  title: string;
}

export interface IncidentAssignee {
  id: string;
  name: string;
  photoUri: string;
  role: string;
}

export interface IncidentsContextType {
  incidents: Incident[];
  isIncidentListLoading: boolean;
  isIncidentPreviewLoading: boolean;
  selectedIncident: Incident | null;
  setIncidents: (incidents: Incident[]) => void;
  setIsIncidentPreviewLoading: (isLoading: boolean) => void;
  setSelectedIncident: (incident: Incident | null) => void;
}

interface IncidentsProviderProps {
  children: ReactNode;
}

export interface IIncidentTask {
  assignee: IncidentAssignee;
  comments: IncidentComment[];
  id: string;
  isDone: boolean;
  items: IncidentTaskItem[];
  priority: number;
  taskName: string;
}

export interface IncidentComment {
  content: IncidentTaskItemContent;
  source: IncidentAssignee;
  timestamp: number;
}

export interface IncidentTaskItem {
  content: IncidentTaskItemContent[];
  source: string;
}

interface IncidentTaskItemContent {
  key: string;
  value: string;
  valueType: 'graph' | 'markdown' | 'text';
}

interface IncidentResponseBody {
  incident: Incident;
}

interface IncidentsResponseBody {
  incidents: Incident[];
}

export const IncidentsContext = createContext<IncidentsContextType | null>(
  null
);

export const IncidentsProvider = ({ children }: IncidentsProviderProps) => {
  const [incidents, setIncidents] = useState<Incident[]>([]);
  const [isIncidentListLoading, setIsIncidentListLoading] = useState(false);
  const [isIncidentPreviewLoading, setIsIncidentPreviewLoading] =
    useState(false);
  const [selectedIncident, setSelectedIncident] = useState<Incident | null>(
    null
  );

  const { accessToken, namespaceId } = useAuthContext();
  const { incidentId } = useParams<{ incidentId: string }>();

  useEffect(() => {
    if (!namespaceId || !accessToken) return;

    const fetchIncident = async () => {
      try {
        setIsIncidentPreviewLoading(true);

        const response = await fetch(
          `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/${incidentId}`,
          {
            headers: {
              Authorization: `Bearer ${accessToken}`
            }
          }
        );

        if (!response.ok) throw new Error('Failed to fetch incident');

        const data: IncidentResponseBody = await response.json();

        setSelectedIncident(data.incident);
      } catch (error) {
        if (error instanceof Error) {
          toast.error(error.message);
        } else {
          toast.error('An unexpected error occurred. Please try again later.');
        }
      } finally {
        setIsIncidentPreviewLoading(false);
      }
    };

    if (incidentId) {
      setSelectedIncident(null);
      fetchIncident();
    } else {
      setSelectedIncident(null);
    }
  }, [accessToken, incidentId, namespaceId]);

  useEffect(() => {
    if (!namespaceId) return;

    setIsIncidentListLoading(true);

    const fetchIncidents = async () => {
      if (!namespaceId || !accessToken) return;

      try {
        const response = await fetch(
          `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/incidents`,
          {
            headers: {
              Authorization: `Bearer ${accessToken}`
            }
          }
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
        setIsIncidentListLoading(false);
      }
    };

    fetchIncidents();
  }, [accessToken, namespaceId]);

  const VALUE = {
    incidents,
    isIncidentListLoading,
    isIncidentPreviewLoading,
    selectedIncident,
    setIncidents,
    setIsIncidentPreviewLoading,
    setSelectedIncident
  };

  return (
    <IncidentsContext.Provider value={VALUE}>
      {children}
    </IncidentsContext.Provider>
  );
};
