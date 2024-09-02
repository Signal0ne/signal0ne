import {
  BackStageIcon,
  ElasticSearchIcon,
  JaegerIcon,
  OpenAIIcon,
  PrometheusIcon,
  ScheduledIcon,
  Signal0neLogo,
  SlackIcon,
  WebhookIcon
} from '../components/Icons/Icons';
import { ReactNode } from 'react';

export const checkDisplayScrollOffset = (element: HTMLElement) => {
  if (!element) return false;

  return element.scrollHeight > element.clientHeight;
};

export const getIntegrationIcon = (integrationName: string) => {
  const icons: Record<string, ReactNode> = {
    alertmanager: <PrometheusIcon />,
    backstage: <BackStageIcon />,
    jaeger: <JaegerIcon />,
    openai: <OpenAIIcon />,
    opensearch: <ElasticSearchIcon />,
    prometheus: <PrometheusIcon />,
    scheduled: <ScheduledIcon />,
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
