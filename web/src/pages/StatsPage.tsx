import { useQuery } from "@tanstack/react-query";
import { Card, Col, Row, Select, Space, Statistic, Table, Typography } from "antd";
import { useMemo, useState } from "react";
import * as ledgersApi from "../api/ledgers";
import * as statsApi from "../api/stats";
import { centsToYuanLabel } from "../utils/money";

export function StatsPage() {
  const [ledgerId, setLedgerId] = useState<number | undefined>();
  const [year, setYear] = useState<number | undefined>();

  const { data: ledgers } = useQuery({
    queryKey: ["ledgers"],
    queryFn: () => ledgersApi.listLedgers(),
  });

  const { data: contacts, isLoading: cLoading } = useQuery({
    queryKey: ["stats", "contacts"],
    queryFn: () => statsApi.contactSummaries(),
  });

  const { data: summary, isLoading: sLoading } = useQuery({
    queryKey: ["stats", "summary", ledgerId, year],
    queryFn: () => statsApi.summary({ ledger_id: ledgerId, year }),
  });

  const yearOptions = useMemo(() => {
    const y = new Date().getFullYear();
    return Array.from({ length: 10 }, (_, i) => y - i);
  }, []);

  return (
    <div>
      <Typography.Title level={4}>汇总</Typography.Title>
      <Space style={{ marginBottom: 16 }} wrap>
        <Select
          allowClear
          placeholder="账本（可选）"
          style={{ width: 180 }}
          value={ledgerId}
          onChange={(v) => setLedgerId(v)}
          options={ledgers?.map((l) => ({ value: l.id, label: l.name }))}
        />
        <Select
          allowClear
          placeholder="年份（可选）"
          style={{ width: 140 }}
          value={year}
          onChange={(v) => setYear(v)}
          options={yearOptions.map((y) => ({ value: y, label: String(y) }))}
        />
      </Space>
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={8}>
          <Card loading={sLoading}>
            <Statistic title="收礼合计(元)" value={centsToYuanLabel(summary?.total_receive_cents ?? 0)} />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card loading={sLoading}>
            <Statistic title="送礼合计(元)" value={centsToYuanLabel(summary?.total_give_cents ?? 0)} />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card loading={sLoading}>
            <Statistic title="差额(元)" value={centsToYuanLabel(summary?.balance_cents ?? 0)} />
          </Card>
        </Col>
      </Row>
      <Typography.Title level={4}>按月</Typography.Title>
      <Table
        rowKey="year_month"
        loading={sLoading}
        dataSource={summary?.monthly}
        pagination={false}
        columns={[
          { title: "月份", dataIndex: "year_month" },
          {
            title: "收礼(元)",
            dataIndex: "receive_cents",
            render: (c: number) => centsToYuanLabel(c),
          },
          {
            title: "送礼(元)",
            dataIndex: "give_cents",
            render: (c: number) => centsToYuanLabel(c),
          },
        ]}
      />
      <Typography.Title level={4} style={{ marginTop: 32 }}>
        按联系人
      </Typography.Title>
      <Table
        rowKey="contact_id"
        loading={cLoading}
        dataSource={contacts}
        pagination={{ pageSize: 20 }}
        columns={[
          { title: "联系人 ID", dataIndex: "contact_id", width: 100 },
          {
            title: "收礼(元)",
            dataIndex: "total_receive_cents",
            render: (c: number) => centsToYuanLabel(c),
          },
          {
            title: "送礼(元)",
            dataIndex: "total_give_cents",
            render: (c: number) => centsToYuanLabel(c),
          },
          {
            title: "差额(元)",
            dataIndex: "balance_cents",
            render: (c: number) => centsToYuanLabel(c),
          },
          {
            title: "最近一笔",
            dataIndex: "last_occurred_on",
            render: (v: string) => v?.slice(0, 10) ?? "—",
          },
        ]}
      />
    </div>
  );
}
