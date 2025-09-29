import { httpRequest } from './httpClient';

export interface Post {
  id: string;
  authorId: string;
  authorName: string;
  content: string;
  createdAt: string;
}

export interface CreatePostPayload {
  content: string;
}

export const postsApi = {
  getFeed: (token: string | null) =>
    httpRequest<Post[]>('/posts', {
      method: 'GET',
      authToken: token ?? undefined,
    }),
  create: (token: string, payload: CreatePostPayload) =>
    httpRequest<Post>('/posts', {
      method: 'POST',
      authToken: token,
      body: payload,
    }),
};
