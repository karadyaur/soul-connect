import { httpRequest } from './httpClient';

export interface LoginPayload {
  email: string;
  password: string;
}

export interface RegisterPayload {
  name: string;
  email: string;
  password: string;
}

export interface AuthResponse<TUser = AuthUser> {
  accessToken: string;
  refreshToken?: string;
  user: TUser;
}

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  bio?: string;
  avatarUrl?: string;
}

export const authApi = {
  login: (payload: LoginPayload) =>
    httpRequest<AuthResponse>('/auth/login', {
      method: 'POST',
      body: payload,
    }),
  register: (payload: RegisterPayload) =>
    httpRequest<AuthResponse>('/auth/register', {
      method: 'POST',
      body: payload,
    }),
  refresh: (token: string) =>
    httpRequest<AuthResponse>('/auth/refresh', {
      method: 'POST',
      authToken: token,
    }),
};
