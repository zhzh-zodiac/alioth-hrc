import { api } from "./client";

export async function listReminders(): Promise<{ items: unknown[]; hint: string }> {
  const { data } = await api.get<{ items: unknown[]; hint: string }>("/reminders");
  return data;
}
