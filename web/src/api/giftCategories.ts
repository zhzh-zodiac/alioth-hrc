import type { GiftCategory } from "../types/api";
import { api } from "./client";

export async function listGiftCategories(): Promise<GiftCategory[]> {
  const { data } = await api.get<{ items: GiftCategory[] }>("/gift-categories");
  return data.items;
}

export async function createGiftCategory(body: { name: string }): Promise<GiftCategory> {
  const { data } = await api.post<GiftCategory>("/gift-categories", body);
  return data;
}

export async function deleteGiftCategory(id: number): Promise<void> {
  await api.delete(`/gift-categories/${id}`);
}
