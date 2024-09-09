import { Integration } from '../contexts/IntegrationsProvider/IntegrationsProvider';

export const DUMMY_AVAILABLE_INTEGRATIONS: Integration[] = [
  {
    config: {
      apiKey: 'string',
      url: 'url'
    },
    imageUri: '../logos/backstage.svg',
    name: 'Backstage',
    type: 'backstage'
  },
  {
    config: {
      url: 'url'
    },
    imageUri: '../logos/jaeger.svg',
    name: 'Jaeger',
    type: 'jaeger'
  },
  {
    config: {
      url: 'url'
    },
    imageUri: '../logos/alertmanager.svg',
    name: 'Alertmanager',
    type: 'alertmanager'
  },
  {
    config: null,
    imageUri: '../logos/signal0ne.svg',
    name: 'Signal0ne',
    type: 'signal0ne'
  },
  {
    config: {
      url: 'url',
      workspaceId: 'string'
    },
    imageUri: '../logos/slack.svg',
    name: 'Slack',
    type: 'slack'
  }
];
