import {
  BackStageIcon,
  JaegerIcon,
  PrometheusIcon,
  ScheduledIcon,
  Signal0neLogo,
  SlackIcon,
  WebhookIcon
} from '../components/Icons/Icons';
import { ReactNode } from 'react';

export const getIntegrationIcon = (integrationName: string) => {
  const icons: Record<string, ReactNode> = {
    backstage: <BackStageIcon />,
    jaeger: <JaegerIcon />,
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
