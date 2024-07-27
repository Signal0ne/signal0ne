import AlertsPage from '../pages/AlertsPage/AlertsPage';
import WorkflowsPage from '../pages/WorkflowsPage/WorkflowsPage';

export const ROUTES = [
  {
    Component: AlertsPage,
    isDisabled: true,
    path: '/alerts',
    title: 'Alerts'
  },
  {
    Component: WorkflowsPage,
    isDisabled: false,
    path: '/',
    title: 'Workflows'
  }
];
