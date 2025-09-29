const API_BASE_URL = process.env.REACT_APP_API_URL ?? 'http://localhost:8080';

export interface HttpRequestOptions
  extends Omit<RequestInit, 'body' | 'headers'> {
  authToken?: string | null;
  parseJson?: boolean;
  body?: RequestInit['body'] | null | object;
  headers?: HeadersInit;
}

export interface HttpErrorPayload {
  message?: string;
  [key: string]: unknown;
}

export class HttpError extends Error {
  public status: number;
  public payload?: HttpErrorPayload;

  constructor(status: number, message: string, payload?: HttpErrorPayload) {
    super(message);
    this.status = status;
    this.payload = payload;
  }
}

const resolveUrl = (path: string): string => {
  if (path.startsWith('http://') || path.startsWith('https://')) {
    return path;
  }

  if (!path.startsWith('/')) {
    return `${API_BASE_URL}/${path}`;
  }

  return `${API_BASE_URL}${path}`;
};

export async function httpRequest<T>(
  path: string,
  options: HttpRequestOptions = {}
): Promise<T> {
  const { authToken, parseJson = true, headers, body, ...rest } = options;
  const requestHeaders = new Headers(headers);

  if (authToken) {
    requestHeaders.set('Authorization', `Bearer ${authToken}`);
  }

  const hasJsonBody = body && !(body instanceof FormData) && typeof body === 'object';
  if (hasJsonBody) {
    requestHeaders.set('Content-Type', 'application/json');
  }

  const response = await fetch(resolveUrl(path), {
    ...rest,
    body: hasJsonBody ? JSON.stringify(body) : body ?? undefined,
    headers: requestHeaders,
    mode: 'cors',
    credentials: 'include',
  });

  if (!response.ok) {
    let payload: HttpErrorPayload | undefined;
    try {
      payload = await response.json();
    } catch (error) {
      // Ignore JSON parsing issues for error payloads.
    }

    const message =
      payload?.message ?? response.statusText ?? 'Unexpected server error';
    throw new HttpError(response.status, message, payload);
  }

  if (response.status === 204 || !parseJson) {
    return undefined as unknown as T;
  }

  return (await response.json()) as T;
}
