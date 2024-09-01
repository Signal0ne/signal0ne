import AlertsPage from '../pages/AlertsPage/AlertsPage';
import IntegrationsPage from '../pages/IntegrationsPage/IntegrationsPage';
import WorkflowsPage from '../pages/WorkflowsPage/WorkflowsPage';

export const ROUTES = [
  {
    Component: AlertsPage,
    isDisabled: true,
    path: '/alerts',
    redirectTo: '/login',
    title: 'Alerts',
    unAuthed: false
  },
  {
    Component: IntegrationsPage,
    path: '/integrations',
    title: 'Integrations',
    redirectTo: '/login',
    unAuthed: false
  },
  {
    Component: WorkflowsPage,
    isDisabled: false,
    path: '/',
    redirectTo: '/login',
    title: 'Workflows',
    unAuthed: false
  }
];
