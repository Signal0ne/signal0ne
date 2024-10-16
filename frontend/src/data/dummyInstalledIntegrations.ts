import type { Integration } from '../contexts/IntegrationsProvider/IntegrationsProvider';

export const DUMMY_INSTALLED_INTEGRATIONS: Integration[] = [
  {
    config: {
      apiKey: 'Some API key',
      url: 'https://test.com:8080'
    },
    id: '123',
    imageUri: 'someUri',
    name: 'backstage_prod',
    type: 'backstage'
  },
  {
    config: {
      url: 'https://test.com:8080'
    },
    id: '345',
    imageUri: 'someUri',
    name: 'jaeger_prod',
    type: 'jaeger'
  },
  {
    config: {
      url: 'https://test.com:8080'
    },
    id: '987',
    imageUri: 'someUri',
    name: 'slack_prod',
    type: 'slack'
  }
];
