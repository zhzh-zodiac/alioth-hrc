import { Button, Card, Form, Input, Typography, message } from "antd";
import { Link, useNavigate, useLocation } from "react-router-dom";
import { login } from "../api/auth";

export function LoginPage() {
  const nav = useNavigate();
  const loc = useLocation() as { state?: { from?: string } };
  const from = loc.state?.from ?? "/records";

  async function onFinish(values: { email: string; password: string }) {
    try {
      await login(values.email, values.password);
      message.success("登录成功");
      nav(from, { replace: true });
    } catch {
      message.error("登录失败，请检查邮箱与密码");
    }
  }

  return (
    <div style={{ maxWidth: 400, margin: "80px auto" }}>
      <Typography.Title level={3} style={{ textAlign: "center" }}>
        登录
      </Typography.Title>
      <Card>
        <Form layout="vertical" onFinish={onFinish}>
          <Form.Item name="email" label="邮箱" rules={[{ required: true, type: "email" }]}>
            <Input autoComplete="email" />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true, min: 8 }]}>
            <Input.Password autoComplete="current-password" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block>
            登录
          </Button>
        </Form>
        <div style={{ marginTop: 16, textAlign: "center" }}>
          <Link to="/register">没有账号？注册</Link>
        </div>
      </Card>
    </div>
  );
}
