import {
  BackStageIcon,
  JaegerIcon,
  PrometheusIcon,
  Signal0neLogo,
  SlackIcon
} from '../components/Icons/Icons';
import { ReactNode } from 'react';

export const getIntegrationIcon = (integrationName: string) => {
  const icons: Record<string, ReactNode> = {
    backstage: <BackStageIcon />,
    jaeger: <JaegerIcon />,
    prometheus: <PrometheusIcon />,
    signal0ne: <Signal0neLogo />,
    slack: <SlackIcon />
  };

  return icons[integrationName] || null;
};

export const handleKeyDown =
  (callback: (...arg: unknown[]) => void, disabled?: boolean) =>
  (e: React.KeyboardEvent) => {
    e.key === ' ' && e.preventDefault();

    if (['Enter', ' '].includes(e.key) && !disabled) callback(e);
  };
