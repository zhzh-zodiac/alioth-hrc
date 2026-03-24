import { useQuery } from "@tanstack/react-query";
import { Alert, Spin, Typography } from "antd";
import * as api from "../api/reminders";

export function RemindersPage() {
  const { data, isLoading } = useQuery({
    queryKey: ["reminders"],
    queryFn: () => api.listReminders(),
  });

  return (
    <div>
      <Typography.Title level={4}>提醒</Typography.Title>
      <Spin spinning={isLoading}>
        <Alert type="info" showIcon message={data?.hint ?? "加载中…"} style={{ marginTop: 16 }} />
      </Spin>
    </div>
  );
}
