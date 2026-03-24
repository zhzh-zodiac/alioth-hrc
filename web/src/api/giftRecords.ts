import type { GiftDirection, GiftRecord, Paginated } from "../types/api";
import { api } from "./client";

export interface GiftRecordListParams {
  contact_id?: number;
  ledger_id?: number;
  category_id?: number;
  direction?: GiftDirection | "";
  from_date?: string;
  to_date?: string;
  page?: number;
  page_size?: number;
}

export async function listGiftRecords(params: GiftRecordListParams): Promise<Paginated<GiftRecord>> {
  const { data } = await api.get<{
    items: GiftRecord[];
    total: number;
    page: number;
    page_size: number;
  }>("/gift-records", { params });
  return {
    items: data.items,
    total: data.total,
    page: data.page,
    page_size: data.page_size,
  };
}

export async function getGiftRecord(id: number): Promise<GiftRecord> {
  const { data } = await api.get<GiftRecord>(`/gift-records/${id}`);
  return data;
}

export async function createGiftRecord(body: {
  ledger_id: number;
  contact_id: number;
  category_id: number;
  direction: GiftDirection;
  amount_cents: number;
  occurred_on: string;
  note?: string;
}): Promise<GiftRecord> {
  const { data } = await api.post<GiftRecord>("/gift-records", body);
  return data;
}

export async function updateGiftRecord(
  id: number,
  body: {
    ledger_id: number;
    contact_id: number;
    category_id: number;
    direction: GiftDirection;
    amount_cents: number;
    occurred_on: string;
    note?: string;
  },
): Promise<GiftRecord> {
  const { data } = await api.put<GiftRecord>(`/gift-records/${id}`, body);
  return data;
}

export async function deleteGiftRecord(id: number): Promise<void> {
  await api.delete(`/gift-records/${id}`);
}

/** CSV 导出（走 axios，可复用 401 刷新） */
export async function exportGiftRecordsCsv(params: GiftRecordListParams): Promise<Blob> {
  const { data } = await api.get<Blob>("/gift-records/export.csv", {
    params,
    responseType: "blob",
  });
  return data;
}
