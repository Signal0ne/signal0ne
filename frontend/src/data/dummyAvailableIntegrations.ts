export interface AvailableIntegration {
  config: Record<string, unknown> | null;
  displayName: string;
  imageUri: string;
  typeName: string;
}

export const DUMMY_AVAILABLE_INTEGRATIONS: AvailableIntegration[] = [
  {
    config: {
      apiKey: 'string',
      host: 'string',
      port: 'string'
    },
    displayName: 'Backstage',
    imageUri: '../logos/backstage.svg',
    typeName: 'backstage'
  },
  {
    config: {
      host: 'string',
      port: 'string'
    },
    displayName: 'Jaeger',
    imageUri: '../logos/jaeger.svg',
    typeName: 'jaeger'
  },
  {
    config: {
      host: 'string',
      port: 'string'
    },
    displayName: 'Alertmanager',
    imageUri: '../logos/alertmanager.svg',
    typeName: 'alertmanager'
  },
  {
    config: null,
    displayName: 'Signal0ne',
    imageUri: '../logos/signal0ne.svg',
    typeName: 'signal0ne'
  },
  {
    config: {
      host: 'string',
      port: 'string',
      workspaceId: 'string'
    },
    displayName: 'Slack',
    imageUri: '../logos/slack.svg',
    typeName: 'slack'
  }
];
