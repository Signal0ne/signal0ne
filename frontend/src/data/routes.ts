import AlertsPage from '../pages/AlertsPage/AlertsPage';
import IntegrationsPage from '../pages/IntegrationsPage/IntegrationsPage';
import WorkflowsPage from '../pages/WorkflowsPage/WorkflowsPage';

export const ROUTES = [
  {
    Component: AlertsPage,
    isDisabled: true,
    path: '/alerts',
    title: 'Alerts'
  },
  {
    Component: IntegrationsPage,
    path: '/integrations',
    title: 'Integrations'
  },
  {
    Component: WorkflowsPage,
    isDisabled: false,
    path: '/',
    title: 'Workflows'
  }
];
