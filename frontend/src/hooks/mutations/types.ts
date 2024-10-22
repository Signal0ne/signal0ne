import type { User } from '../../contexts/AuthProvider/AuthProvider';

export interface AuthPayload {
  password: string;
  username: string;
}

export interface AuthResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
}

export interface ConfigData {
  [key: string]: string;
}

export interface NewIntegrationPayload {
  config: ConfigData;
  name: string;
  type: string;
}
