/** 元（界面输入）转分（接口） */
export function yuanToCents(yuan: number): number {
  return Math.round(yuan * 100);
}

export function centsToYuanLabel(cents: number): string {
  return (cents / 100).toFixed(2);
}
