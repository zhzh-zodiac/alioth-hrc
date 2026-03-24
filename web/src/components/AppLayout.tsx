import { Layout, Menu, Button, theme } from "antd";
import {
  UserOutlined,
  BookOutlined,
  TagsOutlined,
  TransactionOutlined,
  BarChartOutlined,
  BellOutlined,
  LogoutOutlined,
} from "@ant-design/icons";
import { Link, Outlet, useLocation, useNavigate } from "react-router-dom";
import { logout } from "../api/auth";
import { getStoredEmail } from "../api/client";

const { Header, Sider, Content } = Layout;

const items = [
  { key: "/records", icon: <TransactionOutlined />, label: <Link to="/records">礼金流水</Link> },
  { key: "/contacts", icon: <UserOutlined />, label: <Link to="/contacts">联系人</Link> },
  { key: "/ledgers", icon: <BookOutlined />, label: <Link to="/ledgers">账本</Link> },
  { key: "/categories", icon: <TagsOutlined />, label: <Link to="/categories">礼金分类</Link> },
  { key: "/stats", icon: <BarChartOutlined />, label: <Link to="/stats">统计</Link> },
  { key: "/reminders", icon: <BellOutlined />, label: <Link to="/reminders">提醒</Link> },
];

export function AppLayout() {
  const loc = useLocation();
  const nav = useNavigate();
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  async function handleLogout() {
    await logout();
    nav("/login");
  }

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Sider breakpoint="lg" collapsedWidth={0}>
        <div style={{ padding: 16, color: "#fff", fontWeight: 600 }}>人情往来</div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[loc.pathname]}
          items={items}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            padding: "0 24px",
            background: colorBgContainer,
            display: "flex",
            alignItems: "center",
            justifyContent: "flex-end",
            gap: 16,
          }}
        >
          <span style={{ color: "#666" }}>{getStoredEmail() ?? ""}</span>
          <Button onClick={handleLogout} icon={<LogoutOutlined />}>
            退出
          </Button>
        </Header>
        <Content style={{ margin: 24 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
