import axios from "axios";
import { API_BASE } from "../config";
import type { TokenPair } from "../types/api";

const ACCESS_KEY = "hrc_access_token";
const REFRESH_KEY = "hrc_refresh_token";
const EMAIL_KEY = "hrc_user_email";

export function getAccessToken(): string | null {
  return localStorage.getItem(ACCESS_KEY);
}

export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_KEY);
}

export function getStoredEmail(): string | null {
  return localStorage.getItem(EMAIL_KEY);
}

export function setTokens(pair: TokenPair): void {
  localStorage.setItem(ACCESS_KEY, pair.access_token);
  localStorage.setItem(REFRESH_KEY, pair.refresh_token);
}

export function setUserEmail(email: string): void {
  localStorage.setItem(EMAIL_KEY, email);
}

export function clearSession(): void {
  localStorage.removeItem(ACCESS_KEY);
  localStorage.removeItem(REFRESH_KEY);
  localStorage.removeItem(EMAIL_KEY);
}

export const api = axios.create({
  baseURL: API_BASE,
  headers: { "Content-Type": "application/json" },
});

api.interceptors.request.use((config) => {
  const t = getAccessToken();
  if (t) {
    config.headers.Authorization = `Bearer ${t}`;
  }
  return config;
});

let refreshPromise: Promise<string> | null = null;

async function refreshAccessToken(): Promise<string> {
  if (refreshPromise) return refreshPromise;
  refreshPromise = (async () => {
    const rt = getRefreshToken();
    if (!rt) throw new Error("no refresh token");
    const { data } = await axios.post<TokenPair>(
      `${API_BASE}/auth/refresh`,
      { refresh_token: rt },
      { headers: { "Content-Type": "application/json" } },
    );
    setTokens(data);
    return data.access_token;
  })()
    .finally(() => {
      refreshPromise = null;
    });
  return refreshPromise;
}

api.interceptors.response.use(
  (res) => res,
  async (err) => {
    const status = err.response?.status;
    const cfg = err.config as typeof err.config & { _retry?: boolean; url?: string };
    if (!cfg || status !== 401) return Promise.reject(err);

    const url = cfg.url ?? "";
    if (url.includes("/auth/login") || url.includes("/auth/register") || url.includes("/auth/refresh")) {
      return Promise.reject(err);
    }
    if (cfg._retry) {
      clearSession();
      window.location.href = "/login";
      return Promise.reject(err);
    }
    cfg._retry = true;
    try {
      await refreshAccessToken();
      const t = getAccessToken();
      if (t) cfg.headers = { ...cfg.headers, Authorization: `Bearer ${t}` };
      return api.request(cfg);
    } catch {
      clearSession();
      window.location.href = "/login";
      return Promise.reject(err);
    }
  },
);
