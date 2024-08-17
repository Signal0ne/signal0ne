import AlertsPage from '../pages/AlertsPage/AlertsPage';
import IntegrationsPage from '../pages/IntegrationsPage/IntegrationsPage';
import LoginPage from '../pages/LoginPage/LoginPage';
import SignUpPage from '../pages/SignUpPage/SignUpPage';
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
  },
  {
    Component: LoginPage,
    isDisabled: true,
    path: '/login',
    redirectTo: '/',
    title: 'Sign In',
    unAuthed: true
  },
  {
    Component: SignUpPage,
    path: '/register',
    redirectTo: '/',
    title: 'Sign Up',
    unAuthed: true
  }
];
