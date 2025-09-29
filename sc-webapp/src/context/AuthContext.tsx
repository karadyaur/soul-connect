import React, {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react';
import { authApi, type AuthResponse, type AuthUser } from '../api/auth';
import { usersApi } from '../api/users';

export interface AuthState {
  user: AuthUser | null;
  token: string | null;
  isLoading: boolean;
}

export interface AuthContextValue extends AuthState {
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

const TOKEN_STORAGE_KEY = 'sc-webapp/token';

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

const persistAuth = (auth: Pick<AuthResponse, 'accessToken' | 'refreshToken'>) => {
  localStorage.setItem(
    TOKEN_STORAGE_KEY,
    JSON.stringify({
      accessToken: auth.accessToken,
      refreshToken: auth.refreshToken,
    })
  );
};

const restoreToken = (): { accessToken: string; refreshToken?: string } | null => {
  try {
    const serialized = localStorage.getItem(TOKEN_STORAGE_KEY);
    return serialized ? JSON.parse(serialized) : null;
  } catch (error) {
    console.warn('Unable to restore auth token from storage', error);
    return null;
  }
};

const clearPersistedAuth = () => {
  localStorage.removeItem(TOKEN_STORAGE_KEY);
};

export const AuthProvider: React.FC<React.PropsWithChildren> = ({ children }) => {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const bootstrap = useCallback(async () => {
    const stored = restoreToken();
    if (!stored?.accessToken) {
      setIsLoading(false);
      return;
    }

    setToken(stored.accessToken);
    try {
      const profile = await usersApi.getCurrent(stored.accessToken);
      setUser(profile);
    } catch (error) {
      console.error('Failed to load current user profile', error);
      clearPersistedAuth();
      setToken(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void bootstrap();
  }, [bootstrap]);

  const handleAuthSuccess = useCallback(async (response: AuthResponse) => {
    persistAuth(response);
    setToken(response.accessToken);
    setUser(response.user);
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const response = await authApi.login({ email, password });
    await handleAuthSuccess(response);
  }, [handleAuthSuccess]);

  const register = useCallback(
    async (name: string, email: string, password: string) => {
      const response = await authApi.register({ name, email, password });
      await handleAuthSuccess(response);
    },
    [handleAuthSuccess]
  );

  const refreshUser = useCallback(async () => {
    if (!token) {
      return;
    }

    try {
      const profile = await usersApi.getCurrent(token);
      setUser(profile);
    } catch (error) {
      console.error('Failed to refresh profile', error);
    }
  }, [token]);

  const logout = useCallback(() => {
    clearPersistedAuth();
    setToken(null);
    setUser(null);
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({ user, token, isLoading, login, register, logout, refreshUser }),
    [user, token, isLoading, login, register, logout, refreshUser]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuthContext = (): AuthContextValue => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuthContext must be used within AuthProvider');
  }

  return context;
};
