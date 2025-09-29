import { httpRequest } from './httpClient';
import type { AuthUser } from './auth';

export interface UpdateProfilePayload {
  name?: string;
  bio?: string;
}

export interface Profile extends AuthUser {
  bio?: string;
  avatarUrl?: string;
}

export const usersApi = {
  getCurrent: (token: string) =>
    httpRequest<Profile>('/users/me', {
      method: 'GET',
      authToken: token,
    }),
  updateProfile: (token: string, payload: UpdateProfilePayload) =>
    httpRequest<Profile>('/users/me', {
      method: 'PUT',
      authToken: token,
      body: payload,
    }),
};
