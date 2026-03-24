import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  Button,
  DatePicker,
  Form,
  Input,
  InputNumber,
  Modal,
  Select,
  Space,
  Table,
  message,
} from "antd";
import type { Dayjs } from "dayjs";
import dayjs from "dayjs";
import { useMemo, useState } from "react";
import * as contactsApi from "../api/contacts";
import * as giftCatApi from "../api/giftCategories";
import * as giftRecApi from "../api/giftRecords";
import * as ledgersApi from "../api/ledgers";
import type { GiftDirection, GiftRecord } from "../types/api";
import { centsToYuanLabel, yuanToCents } from "../utils/money";

export function RecordsPage() {
  const qc = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [contactId, setContactId] = useState<number | undefined>();
  const [ledgerId, setLedgerId] = useState<number | undefined>();
  const [categoryId, setCategoryId] = useState<number | undefined>();
  const [direction, setDirection] = useState<GiftDirection | "">("");
  const [range, setRange] = useState<[Dayjs | null, Dayjs | null] | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<GiftRecord | null>(null);
  const [form] = Form.useForm<{
    ledger_id: number;
    contact_id: number;
    category_id: number;
    direction: GiftDirection;
    amount_yuan: number;
    occurred_on: Dayjs;
    note?: string;
  }>();

  const listParams = useMemo(
    () => ({
      page,
      page_size: pageSize,
      contact_id: contactId,
      ledger_id: ledgerId,
      category_id: categoryId,
      direction: direction || undefined,
      from_date: range?.[0]?.format("YYYY-MM-DD"),
      to_date: range?.[1]?.format("YYYY-MM-DD"),
    }),
    [page, pageSize, contactId, ledgerId, categoryId, direction, range],
  );

  const { data: contacts } = useQuery({
    queryKey: ["contacts", "all"],
    queryFn: () => contactsApi.listContacts({ page: 1, page_size: 500 }),
  });
  const { data: ledgers } = useQuery({
    queryKey: ["ledgers"],
    queryFn: () => ledgersApi.listLedgers(),
  });
  const { data: categories } = useQuery({
    queryKey: ["gift-categories"],
    queryFn: () => giftCatApi.listGiftCategories(),
  });

  const { data, isLoading } = useQuery({
    queryKey: ["gift-records", listParams],
    queryFn: () => giftRecApi.listGiftRecords(listParams),
  });

  const saveMut = useMutation({
    mutationFn: async (values: {
      ledger_id: number;
      contact_id: number;
      category_id: number;
      direction: GiftDirection;
      amount_yuan: number;
      occurred_on: Dayjs;
      note?: string;
    }) => {
      const body = {
        ledger_id: values.ledger_id,
        contact_id: values.contact_id,
        category_id: values.category_id,
        direction: values.direction,
        amount_cents: yuanToCents(values.amount_yuan),
        occurred_on: values.occurred_on.format("YYYY-MM-DD"),
        note: values.note ?? "",
      };
      if (editing) return giftRecApi.updateGiftRecord(editing.id, body);
      return giftRecApi.createGiftRecord(body);
    },
    onSuccess: () => {
      message.success("已保存");
      setModalOpen(false);
      setEditing(null);
      form.resetFields();
      void qc.invalidateQueries({ queryKey: ["gift-records"] });
      void qc.invalidateQueries({ queryKey: ["stats"] });
    },
    onError: (e: Error & { response?: { data?: { error?: string } } }) => {
      message.error(e.response?.data?.error ?? "保存失败");
    },
  });

  const delMut = useMutation({
    mutationFn: (id: number) => giftRecApi.deleteGiftRecord(id),
    onSuccess: () => {
      message.success("已删除");
      void qc.invalidateQueries({ queryKey: ["gift-records"] });
      void qc.invalidateQueries({ queryKey: ["stats"] });
    },
    onError: () => message.error("删除失败"),
  });

  const exportMut = useMutation({
    mutationFn: () =>
      giftRecApi.exportGiftRecordsCsv({
        contact_id: contactId,
        ledger_id: ledgerId,
        category_id: categoryId,
        direction: direction || undefined,
        from_date: range?.[0]?.format("YYYY-MM-DD"),
        to_date: range?.[1]?.format("YYYY-MM-DD"),
      }),
    onSuccess: (blob) => {
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = "gift_records.csv";
      a.click();
      URL.revokeObjectURL(url);
      message.success("已开始下载");
    },
    onError: () => message.error("导出失败"),
  });

  function openCreate() {
    setEditing(null);
    form.resetFields();
    const defaultLedger = ledgers?.find((l) => l.is_default);
    form.setFieldsValue({
      ledger_id: defaultLedger?.id,
      direction: "give",
      occurred_on: dayjs(),
      amount_yuan: 1,
    });
    setModalOpen(true);
  }

  function openEdit(row: GiftRecord) {
    setEditing(row);
    form.setFieldsValue({
      ledger_id: row.ledger_id,
      contact_id: row.contact_id,
      category_id: row.category_id,
      direction: row.direction,
      amount_yuan: row.amount_cents / 100,
      occurred_on: dayjs(row.occurred_on),
      note: row.note,
    });
    setModalOpen(true);
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }} wrap>
        <Select
          allowClear
          placeholder="联系人"
          style={{ width: 160 }}
          value={contactId}
          onChange={(v) => {
            setContactId(v);
            setPage(1);
          }}
          options={contacts?.items.map((c) => ({ value: c.id, label: c.name }))}
        />
        <Select
          allowClear
          placeholder="账本"
          style={{ width: 140 }}
          value={ledgerId}
          onChange={(v) => {
            setLedgerId(v);
            setPage(1);
          }}
          options={ledgers?.map((l) => ({ value: l.id, label: l.name }))}
        />
        <Select
          allowClear
          placeholder="分类"
          style={{ width: 140 }}
          value={categoryId}
          onChange={(v) => {
            setCategoryId(v);
            setPage(1);
          }}
          options={categories?.map((c) => ({ value: c.id, label: c.name }))}
        />
        <Select
          allowClear
          placeholder="方向"
          style={{ width: 120 }}
          value={direction || undefined}
          onChange={(v) => {
            setDirection(v ?? "");
            setPage(1);
          }}
          options={[
            { value: "give", label: "送出" },
            { value: "receive", label: "收到" },
          ]}
        />
        <DatePicker.RangePicker
          value={range}
          onChange={(v) => {
            setRange(v);
            setPage(1);
          }}
        />
        <Button type="primary" onClick={openCreate}>
          新建流水
        </Button>
        <Button loading={exportMut.isPending} onClick={() => exportMut.mutate()}>
          导出 CSV
        </Button>
      </Space>
      <Table<GiftRecord>
        rowKey="id"
        loading={isLoading}
        dataSource={data?.items}
        scroll={{ x: 900 }}
        pagination={{
          current: page,
          pageSize,
          total: data?.total ?? 0,
          showSizeChanger: true,
          onChange: (p, ps) => {
            setPage(p);
            setPageSize(ps);
          },
        }}
        columns={[
          {
            title: "日期",
            dataIndex: "occurred_on",
            width: 120,
            render: (v: string) => v?.slice(0, 10),
          },
          {
            title: "方向",
            dataIndex: "direction",
            width: 80,
            render: (v: GiftDirection) => (v === "give" ? "送出" : "收到"),
          },
          {
            title: "金额(元)",
            dataIndex: "amount_cents",
            width: 100,
            render: (c: number) => centsToYuanLabel(c),
          },
          {
            title: "联系人",
            width: 120,
            render: (_, row) => contacts?.items.find((c) => c.id === row.contact_id)?.name ?? row.contact_id,
          },
          {
            title: "账本",
            width: 120,
            render: (_, row) => ledgers?.find((l) => l.id === row.ledger_id)?.name ?? row.ledger_id,
          },
          {
            title: "分类",
            width: 100,
            render: (_, row) => categories?.find((c) => c.id === row.category_id)?.name ?? row.category_id,
          },
          { title: "备注", dataIndex: "note", ellipsis: true },
          {
            title: "操作",
            width: 140,
            fixed: "right",
            render: (_, row) => (
              <Space>
                <Button type="link" size="small" onClick={() => openEdit(row)}>
                  编辑
                </Button>
                <Button
                  type="link"
                  danger
                  size="small"
                  onClick={() => {
                    Modal.confirm({
                      title: "确认删除？",
                      onOk: () => delMut.mutate(row.id),
                    });
                  }}
                >
                  删除
                </Button>
              </Space>
            ),
          },
        ]}
      />
      <Modal
        title={editing ? "编辑流水" : "新建流水"}
        open={modalOpen}
        width={520}
        onCancel={() => {
          setModalOpen(false);
          setEditing(null);
          form.resetFields();
        }}
        onOk={() => void form.submit()}
        confirmLoading={saveMut.isPending}
      >
        <Form form={form} layout="vertical" onFinish={(v) => saveMut.mutate(v)}>
          <Form.Item name="ledger_id" label="账本" rules={[{ required: true }]}>
            <Select options={ledgers?.map((l) => ({ value: l.id, label: l.name }))} />
          </Form.Item>
          <Form.Item name="contact_id" label="联系人" rules={[{ required: true }]}>
            <Select showSearch optionFilterProp="label" options={contacts?.items.map((c) => ({ value: c.id, label: c.name }))} />
          </Form.Item>
          <Form.Item name="category_id" label="分类" rules={[{ required: true }]}>
            <Select options={categories?.map((c) => ({ value: c.id, label: c.name }))} />
          </Form.Item>
          <Form.Item name="direction" label="方向" rules={[{ required: true }]}>
            <Select
              options={[
                { value: "give", label: "送出" },
                { value: "receive", label: "收到" },
              ]}
            />
          </Form.Item>
          <Form.Item name="amount_yuan" label="金额(元)" rules={[{ required: true }]}>
            <InputNumber min={0.01} step={0.01} style={{ width: "100%" }} />
          </Form.Item>
          <Form.Item name="occurred_on" label="发生日期" rules={[{ required: true }]}>
            <DatePicker style={{ width: "100%" }} />
          </Form.Item>
          <Form.Item name="note" label="备注">
            <Input.TextArea rows={2} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
