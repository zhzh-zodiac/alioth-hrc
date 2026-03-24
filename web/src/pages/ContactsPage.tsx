import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Button, Form, Input, Modal, Space, Table, message } from "antd";
import { useState } from "react";
import * as api from "../api/contacts";
import type { Contact } from "../types/api";

export function ContactsPage() {
  const qc = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [search, setSearch] = useState("");
  const [editing, setEditing] = useState<Contact | null>(null);
  const [creating, setCreating] = useState(false);
  const [form] = Form.useForm();

  const { data, isLoading } = useQuery({
    queryKey: ["contacts", page, pageSize, search],
    queryFn: () => api.listContacts({ page, page_size: pageSize, q: search || undefined }),
  });

  const saveMut = useMutation({
    mutationFn: async (values: { name: string; relation_note?: string; remark?: string }) => {
      if (editing) return api.updateContact(editing.id, values);
      return api.createContact(values);
    },
    onSuccess: () => {
      message.success("已保存");
      setEditing(null);
      setCreating(false);
      form.resetFields();
      void qc.invalidateQueries({ queryKey: ["contacts"] });
    },
    onError: () => message.error("保存失败"),
  });

  const delMut = useMutation({
    mutationFn: (id: number) => api.deleteContact(id),
    onSuccess: () => {
      message.success("已删除");
      void qc.invalidateQueries({ queryKey: ["contacts"] });
    },
    onError: (e: Error) => message.error(e.message || "删除失败"),
  });

  return (
    <div>
      <Space style={{ marginBottom: 16 }} wrap>
        <Input.Search
          placeholder="搜索姓名/关系"
          allowClear
          style={{ width: 240 }}
          onSearch={(v) => {
            setSearch(v);
            setPage(1);
          }}
        />
        <Button
          type="primary"
          onClick={() => {
            setCreating(true);
            setEditing(null);
            form.resetFields();
          }}
        >
          新建联系人
        </Button>
      </Space>
      <Table<Contact>
        rowKey="id"
        loading={isLoading}
        dataSource={data?.items}
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
          { title: "姓名", dataIndex: "name" },
          { title: "关系", dataIndex: "relation_note" },
          { title: "备注", dataIndex: "remark", ellipsis: true },
          {
            title: "操作",
            render: (_, row) => (
              <Space>
                <Button
                  type="link"
                  size="small"
                  onClick={() => {
                    setEditing(row);
                    setCreating(true);
                    form.setFieldsValue(row);
                  }}
                >
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
        title={editing ? "编辑联系人" : "新建联系人"}
        open={creating}
        onCancel={() => {
          setCreating(false);
          setEditing(null);
          form.resetFields();
        }}
        onOk={() => void form.submit()}
        confirmLoading={saveMut.isPending}
      >
        <Form form={form} layout="vertical" onFinish={(v) => saveMut.mutate(v)}>
          <Form.Item name="name" label="姓名" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="relation_note" label="关系">
            <Input />
          </Form.Item>
          <Form.Item name="remark" label="备注">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
