import { IncidentsProvider } from '../../contexts/IncidentsProvider/IncidentsProvider';
import IncidentPreview from '../../components/IncidentPreview/IncidentPreview';
import IncidentsSidebar from '../../components/IncidentsSidebar/IncidentsSidebar';
import './IncidentsPage.scss';

const IncidentsPage = () => {
  return (
    <IncidentsProvider>
      <div className="incidents-page">
        <IncidentsSidebar />
        <IncidentPreview />
      </div>
    </IncidentsProvider>
  );
};

export default IncidentsPage;
