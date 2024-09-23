import { IncidentsContext } from '../contexts/IncidentsProvider/IncidentsProvider';
import { useContext } from 'react';

export const useIncidentsContext = () => {
  const context = useContext(IncidentsContext);

  if (!context) {
    throw new Error(
      'useIncidentsContext must be used within an IncidentsProvider'
    );
  }

  return context;
};
