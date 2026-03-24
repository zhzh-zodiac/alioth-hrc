export interface TokenPair {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
  refresh_expires_in: number;
}

export interface Contact {
  id: number;
  user_id: number;
  name: string;
  relation_note: string;
  remark: string;
  created_at: string;
  updated_at: string;
}

export interface Ledger {
  id: number;
  user_id: number;
  name: string;
  is_default: boolean;
  created_at: string;
  updated_at: string;
}

export interface GiftCategory {
  id: number;
  user_id?: number;
  name: string;
  is_system: boolean;
  created_at: string;
  updated_at: string;
}

export type GiftDirection = "give" | "receive";

export interface GiftRecord {
  id: number;
  user_id: number;
  ledger_id: number;
  contact_id: number;
  category_id: number;
  direction: GiftDirection;
  amount_cents: number;
  occurred_on: string;
  note: string;
  created_at: string;
  updated_at: string;
}

export interface ContactStatRow {
  contact_id: number;
  total_receive_cents: number;
  total_give_cents: number;
  balance_cents: number;
  last_occurred_on?: string;
}

export interface MonthlyStatRow {
  year_month: string;
  receive_cents: number;
  give_cents: number;
}

export interface SummaryResult {
  ledger_id?: number;
  year?: number;
  monthly: MonthlyStatRow[];
  total_receive_cents: number;
  total_give_cents: number;
  balance_cents: number;
}

export interface Paginated<T> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
}
