import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Button, Form, Input, Modal, Space, Table, Tag, message } from "antd";
import { useState } from "react";
import * as api from "../api/giftCategories";
import type { GiftCategory } from "../types/api";

export function CategoriesPage() {
  const qc = useQueryClient();
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm();

  const { data, isLoading } = useQuery({
    queryKey: ["gift-categories"],
    queryFn: () => api.listGiftCategories(),
  });

  const createMut = useMutation({
    mutationFn: (values: { name: string }) => api.createGiftCategory(values),
    onSuccess: () => {
      message.success("已创建");
      setOpen(false);
      form.resetFields();
      void qc.invalidateQueries({ queryKey: ["gift-categories"] });
    },
    onError: () => message.error("创建失败"),
  });

  const delMut = useMutation({
    mutationFn: (id: number) => api.deleteGiftCategory(id),
    onSuccess: () => {
      message.success("已删除");
      void qc.invalidateQueries({ queryKey: ["gift-categories"] });
    },
    onError: (e: Error & { response?: { data?: { error?: string } } }) => {
      message.error(e.response?.data?.error ?? "删除失败");
    },
  });

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button type="primary" onClick={() => setOpen(true)}>
          新建分类
        </Button>
      </Space>
      <Table<GiftCategory>
        rowKey="id"
        loading={isLoading}
        dataSource={data}
        columns={[
          { title: "名称", dataIndex: "name" },
          {
            title: "类型",
            render: (_, row) => (row.is_system ? <Tag>系统</Tag> : <Tag color="green">自定义</Tag>),
          },
          {
            title: "操作",
            render: (_, row) =>
              row.is_system ? (
                "—"
              ) : (
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
              ),
          },
        ]}
      />
      <Modal
        title="新建礼金分类"
        open={open}
        onCancel={() => {
          setOpen(false);
          form.resetFields();
        }}
        onOk={() => void form.submit()}
        confirmLoading={createMut.isPending}
      >
        <Form form={form} layout="vertical" onFinish={(v) => createMut.mutate(v)}>
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
