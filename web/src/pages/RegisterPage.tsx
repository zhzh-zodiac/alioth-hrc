import { Button, Card, Form, Input, Typography, message } from "antd";
import { Link, useNavigate } from "react-router-dom";
import { register } from "../api/auth";

export function RegisterPage() {
  const nav = useNavigate();

  async function onFinish(values: { email: string; password: string; name?: string }) {
    try {
      await register(values);
      message.success("注册成功，请登录");
      nav("/login");
    } catch {
      message.error("注册失败，邮箱可能已被使用");
    }
  }

  return (
    <div style={{ maxWidth: 400, margin: "80px auto" }}>
      <Typography.Title level={3} style={{ textAlign: "center" }}>
        注册
      </Typography.Title>
      <Card>
        <Form layout="vertical" onFinish={onFinish}>
          <Form.Item name="email" label="邮箱" rules={[{ required: true, type: "email" }]}>
            <Input autoComplete="email" />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true, min: 8, max: 72 }]}>
            <Input.Password autoComplete="new-password" />
          </Form.Item>
          <Form.Item name="name" label="称呼（可选）">
            <Input />
          </Form.Item>
          <Button type="primary" htmlType="submit" block>
            注册
          </Button>
        </Form>
        <div style={{ marginTop: 16, textAlign: "center" }}>
          <Link to="/login">已有账号？登录</Link>
        </div>
      </Card>
    </div>
  );
}
