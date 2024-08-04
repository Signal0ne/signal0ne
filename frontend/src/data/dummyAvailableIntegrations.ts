export interface AvailableIntegration {
  icon: string;
  name: string;
}

export const DUMMY_AVAILABLE_INTEGRATIONS: AvailableIntegration[] = [
  {
    icon: 'backstage',
    name: 'Backstage'
  },
  {
    icon: 'jaeger',
    name: 'Jaeger'
  },
  {
    icon: 'prometheus',
    name: 'Prometheus'
  },
  {
    icon: 'signal0ne',
    name: 'Signal0ne'
  },
  {
    icon: 'slack',
    name: 'Slack'
  }
];
