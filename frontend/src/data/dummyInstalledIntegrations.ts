export interface InstalledIntegration {
  icon: string;
  name: string;
}

export const DUMMY_INSTALLED_INTEGRATIONS: InstalledIntegration[] = [
  {
    icon: 'backstage',
    name: 'Backstage'
  },
  {
    icon: 'jaeger',
    name: 'Jaeger'
  },
  {
    icon: 'slack',
    name: 'Slack'
  }
];
