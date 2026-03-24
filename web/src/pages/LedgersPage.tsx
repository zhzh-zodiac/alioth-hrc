import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Button, Form, Input, Modal, Space, Table, Tag, message } from "antd";
import { useState } from "react";
import * as api from "../api/ledgers";
import type { Ledger } from "../types/api";

export function LedgersPage() {
  const qc = useQueryClient();
  const [editing, setEditing] = useState<Ledger | null>(null);
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm();

  const { data, isLoading } = useQuery({
    queryKey: ["ledgers"],
    queryFn: () => api.listLedgers(),
  });

  const saveMut = useMutation({
    mutationFn: async (values: { name: string }) => {
      if (editing) return api.updateLedger(editing.id, values);
      return api.createLedger(values);
    },
    onSuccess: () => {
      message.success("已保存");
      setOpen(false);
      setEditing(null);
      form.resetFields();
      void qc.invalidateQueries({ queryKey: ["ledgers"] });
    },
    onError: (e: Error & { response?: { data?: { error?: string } } }) => {
      message.error(e.response?.data?.error ?? "保存失败");
    },
  });

  const delMut = useMutation({
    mutationFn: (id: number) => api.deleteLedger(id),
    onSuccess: () => {
      message.success("已删除");
      void qc.invalidateQueries({ queryKey: ["ledgers"] });
    },
    onError: (e: Error & { response?: { data?: { error?: string } } }) => {
      message.error(e.response?.data?.error ?? "删除失败");
    },
  });

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button
          type="primary"
          onClick={() => {
            setEditing(null);
            form.resetFields();
            setOpen(true);
          }}
        >
          新建账本
        </Button>
      </Space>
      <Table<Ledger>
        rowKey="id"
        loading={isLoading}
        dataSource={data}
        columns={[
          { title: "名称", dataIndex: "name" },
          {
            title: "默认",
            dataIndex: "is_default",
            render: (v: boolean) => (v ? <Tag color="blue">默认</Tag> : "—"),
          },
          {
            title: "操作",
            render: (_, row) => (
              <Space>
                <Button
                  type="link"
                  size="small"
                  onClick={() => {
                    setEditing(row);
                    form.setFieldsValue({ name: row.name });
                    setOpen(true);
                  }}
                >
                  重命名
                </Button>
                <Button
                  type="link"
                  danger
                  size="small"
                  disabled={row.is_default}
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
        title={editing ? "编辑账本" : "新建账本"}
        open={open}
        onCancel={() => {
          setOpen(false);
          setEditing(null);
          form.resetFields();
        }}
        onOk={() => void form.submit()}
        confirmLoading={saveMut.isPending}
      >
        <Form form={form} layout="vertical" onFinish={(v) => saveMut.mutate(v)}>
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
