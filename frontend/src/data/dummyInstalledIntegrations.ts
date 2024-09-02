export interface InstalledIntegration {
  name: string;
  type: string;
}

export const DUMMY_INSTALLED_INTEGRATIONS: InstalledIntegration[] = [
  {
    name: 'backstage_prod',
    type: 'Backstage'
  },
  {
    name: 'jaeger Prod',
    type: 'Jaeger'
  },
  {
    name: 'slack_prod',
    type: 'Slack'
  }
];
