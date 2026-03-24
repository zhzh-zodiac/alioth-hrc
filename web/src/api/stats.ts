import type { ContactStatRow, SummaryResult } from "../types/api";
import { api } from "./client";

export async function contactSummaries(): Promise<ContactStatRow[]> {
  const { data } = await api.get<{ items: ContactStatRow[] }>("/stats/contacts");
  return data.items;
}

export async function summary(params: { ledger_id?: number; year?: number }): Promise<SummaryResult> {
  const { data } = await api.get<SummaryResult>("/stats/summary", { params });
  return data;
}
