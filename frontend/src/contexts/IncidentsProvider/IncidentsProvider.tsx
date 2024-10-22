import { createContext, ReactNode, useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { useGetIncidentByIdQuery } from '../../hooks/queries/useGetIncidentByIdQuery';
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
  isIncidentPreviewLoading: boolean;
  selectedIncident: Incident | null;
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

export const IncidentsContext = createContext<IncidentsContextType | null>(
  null
);

export const IncidentsProvider = ({ children }: IncidentsProviderProps) => {
  const [selectedIncident, setSelectedIncident] = useState<Incident | null>(
    null
  );

  const { incidentId } = useParams<{ incidentId: string }>();

  const { data, isError, isLoading } = useGetIncidentByIdQuery(incidentId);

  useEffect(() => {
    if (isError) toast.error('Failed to get incident data');
  }, [isError]);

  useEffect(() => {
    if (data) setSelectedIncident(data.incident);
  }, [data]);

  useEffect(() => {
    if (!incidentId) setSelectedIncident(null);
  }, [incidentId]);

  const VALUE = {
    isIncidentPreviewLoading: isLoading,
    selectedIncident,
    setSelectedIncident
  };

  return (
    <IncidentsContext.Provider value={VALUE}>
      {children}
    </IncidentsContext.Provider>
  );
};
