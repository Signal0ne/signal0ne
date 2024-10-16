import type { ComponentType, ReactNode } from 'react';
import AlertsPage from '../pages/AlertsPage/AlertsPage';
import IncidentsPage from '../pages/IncidentsPage/IncidentsPage';
import IntegrationsPage from '../pages/IntegrationsPage/IntegrationsPage';
import LoginPage from '../pages/LoginPage/LoginPage';
import SignUpPage from '../pages/SignUpPage/SignUpPage';
import WorkflowsPage from '../pages/WorkflowsPage/WorkflowsPage';

interface RouteConfig {
  Component: ComponentType<{ children?: ReactNode }>;
  isDisabled?: boolean;
  path: string;
  redirectTo?: string;
  showInNavbar: boolean;
  title: string;
  unAuthed: boolean;
}

export const ROUTES: RouteConfig[] = [
  {
    Component: AlertsPage,
    isDisabled: true,
    path: '/alerts',
    redirectTo: '/login',
    title: 'Alerts',
    showInNavbar: true,
    unAuthed: false
  },
  {
    Component: IncidentsPage,
    path: '/incidents',
    title: 'Incidents',
    redirectTo: '/login',
    showInNavbar: true,
    unAuthed: false
  },
  {
    Component: IncidentsPage,
    path: '/incidents/:incidentId',
    title: 'Incidents',
    redirectTo: '/login',
    showInNavbar: false,
    unAuthed: false
  },
  {
    Component: IntegrationsPage,
    path: '/integrations',
    title: 'Integrations',
    redirectTo: '/login',
    showInNavbar: true,
    unAuthed: false
  },
  {
    Component: WorkflowsPage,
    path: '/',
    redirectTo: '/login',
    showInNavbar: true,
    title: 'Workflows',
    unAuthed: false
  },
  {
    Component: WorkflowsPage,
    path: '/:workflowId',
    redirectTo: '/login',
    showInNavbar: false,
    title: 'Workflows',
    unAuthed: false
  },
  {
    Component: LoginPage,
    path: '/login',
    redirectTo: '/',
    showInNavbar: true,
    title: 'Sign In',
    unAuthed: true
  },
  {
    Component: SignUpPage,
    path: '/register',
    redirectTo: '/',
    showInNavbar: true,
    title: 'Sign Up',
    unAuthed: true
  }
];
