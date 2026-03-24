import type { Ledger } from "../types/api";
import { api } from "./client";

export async function listLedgers(): Promise<Ledger[]> {
  const { data } = await api.get<{ items: Ledger[] }>("/ledgers");
  return data.items;
}

export async function getLedger(id: number): Promise<Ledger> {
  const { data } = await api.get<Ledger>(`/ledgers/${id}`);
  return data;
}

export async function createLedger(body: { name: string }): Promise<Ledger> {
  const { data } = await api.post<Ledger>("/ledgers", body);
  return data;
}

export async function updateLedger(id: number, body: { name: string }): Promise<Ledger> {
  const { data } = await api.put<Ledger>(`/ledgers/${id}`, body);
  return data;
}

export async function deleteLedger(id: number): Promise<void> {
  await api.delete(`/ledgers/${id}`);
}
