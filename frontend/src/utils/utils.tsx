import {
  BackStageIcon,
  ConfluenceIcon,
  GithubIcon,
  JaegerIcon,
  OpenAIIcon,
  OpenSearchIcon,
  PagerDutyIcon,
  PrometheusIcon,
  ScheduledIcon,
  ServiceNowIcon,
  Signal0neLogo,
  SlackIcon,
  WebhookIcon
} from '../components/Icons/Icons';
import { ReactNode } from 'react';

export const checkDisplayScrollOffset = (element: HTMLElement) => {
  if (!element) return false;

  return element.scrollHeight > element.clientHeight;
};

export const getFormattedFormLabel = (fieldLabel: string) => {
  switch (fieldLabel) {
    case 'apiKey':
      return 'API Key';
    case 'url':
      return 'URL';
    case 'workspaceId':
      return 'Workspace ID';
    default:
      return fieldLabel;
  }
};

export const getInputType = (fieldLabel: string) => {
  switch (fieldLabel) {
    case 'apiKey':
      return 'password';
    default:
      return 'text';
  }
};

export const getIntegrationGradientColor = (integrationName: string) => {
  switch (integrationName) {
    case 'alertmanager':
      return 'linear-gradient(45deg, #da4e31 0%, #e77f6a 100%)';
    case 'backstage':
      return 'linear-gradient(45deg, #36baa2 0%, #36baa2 100%)';
    case 'confluence':
      return 'linear-gradient(45deg, #2684ff 0%, #59a0fe 100%)';
    case 'github':
      return 'linear-gradient(90deg, #5a6370 0%, #fafbfc 200%)';
    case 'jaeger':
      return 'linear-gradient(90deg, #fff 0%, #60d0e4 50%, #638b18 75%, #e1caa2 100%)';
    case 'openai':
      return 'linear-gradient(135deg, #10A37F 0%, #CAFEFF 200%)';
    case 'opensearch':
      return 'linear-gradient(45deg, #0073b4 0%, #005EB8 100%)';
    case 'pagerduty':
      return 'linear-gradient(45deg, #06AC38 0%, #00dc42 100%)';
    case 'servicenow':
      return 'linear-gradient(45deg, #62d84e 0%, #48a063 100%)';
    case 'signal0ne':
      return 'linear-gradient(45deg, #fff 0%, #eee 100%)';
    case 'slack':
      return 'linear-gradient(90deg, #e01e5a 0%, #ecb22d 33%,#2fb67c 66%,#36c5f1 100%)';
    default:
      return 'linear-gradient(45deg, #fff 0%, #eee 100%)';
  }
};

export const getIntegrationIcon = (integrationName: string) => {
  const icons: Record<string, ReactNode> = {
    alertmanager: <PrometheusIcon />,
    backstage: <BackStageIcon />,
    confluence: <ConfluenceIcon />,
    github: <GithubIcon />,
    jaeger: <JaegerIcon />,
    openai: <OpenAIIcon />,
    opensearch: <OpenSearchIcon />,
    pagerduty: <PagerDutyIcon />,
    prometheus: <PrometheusIcon />,
    scheduled: <ScheduledIcon />,
    servicenow: <ServiceNowIcon />,
    signal0ne: <Signal0neLogo />,
    slack: <SlackIcon />,
    webhook: <WebhookIcon />
  };

  return icons[integrationName] || null;
};

export const handleKeyDown =
  (callback: (...arg: unknown[]) => void, disabled?: boolean) =>
  (e: React.KeyboardEvent) => {
    e.key === ' ' && e.preventDefault();

    if (['Enter', ' '].includes(e.key) && !disabled) callback(e);
  };
