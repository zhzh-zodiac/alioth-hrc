import type { Contact, Paginated } from "../types/api";
import { api } from "./client";

export async function listContacts(params: {
  q?: string;
  page?: number;
  page_size?: number;
}): Promise<Paginated<Contact>> {
  const { data } = await api.get<{ items: Contact[]; total: number; page: number; page_size: number }>(
    "/contacts",
    { params },
  );
  return {
    items: data.items,
    total: data.total,
    page: data.page,
    page_size: data.page_size,
  };
}

export async function getContact(id: number): Promise<Contact> {
  const { data } = await api.get<Contact>(`/contacts/${id}`);
  return data;
}

export async function createContact(body: {
  name: string;
  relation_note?: string;
  remark?: string;
}): Promise<Contact> {
  const { data } = await api.post<Contact>("/contacts", body);
  return data;
}

export async function updateContact(
  id: number,
  body: { name: string; relation_note?: string; remark?: string },
): Promise<Contact> {
  const { data } = await api.put<Contact>(`/contacts/${id}`, body);
  return data;
}

export async function deleteContact(id: number): Promise<void> {
  await api.delete(`/contacts/${id}`);
}
