import axios from "axios";
import { API_BASE } from "../config";
import type { TokenPair } from "../types/api";
import { api, clearSession, setTokens, setUserEmail } from "./client";

export async function login(email: string, password: string): Promise<TokenPair> {
  const { data } = await axios.post<TokenPair>(`${API_BASE}/auth/login`, {
    email,
    password,
  });
  setTokens(data);
  setUserEmail(email);
  return data;
}

export async function register(body: {
  email: string;
  password: string;
  name?: string;
}): Promise<{ id: number; email: string; name: string }> {
  const { data } = await axios.post(`${API_BASE}/auth/register`, body);
  return data;
}

export async function logout(): Promise<void> {
  const refresh_token = localStorage.getItem("hrc_refresh_token");
  if (refresh_token) {
    try {
      await api.post("/auth/logout", { refresh_token });
    } catch {
      // ignore
    }
  }
  clearSession();
}
