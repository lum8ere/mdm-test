// src/App.tsx
import React, { useEffect, useState } from "react";
import { Button, Layout, Menu, message } from "antd";
import DeviceCard from "./components/DeviceCard";
import HeartbeatLog from "./components/HeartbeatLog";
import DeviceList from "./components/DeviceList";
import LoginPage from "./components/LoginPage";
import { jwtDecode } from "jwt-decode";
import axios from "axios";

const { Header, Content, Footer, Sider } = Layout;

interface JwtPayload {
  role: string;
  username: string;
  exp: number;
}

export const App: React.FC = () => {
  const serverUrl = "http://localhost:4000";
  const [token, setToken] = useState<string | null>(
    localStorage.getItem("token")
  );

  // При монтировании устанавливаем заголовок Authorization, если токен уже есть
  useEffect(() => {
    if (token) {
      axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;
    } else {
      delete axios.defaults.headers.common["Authorization"];
    }
  }, [token]);

  const handleLogin = (token: string) => {
    localStorage.setItem("token", token);
    setToken(token);
    // Настроим глобальный заголовок для axios:
    axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;
  };

  const handleLogout = async () => {
    try {
      message.success("Вы вышли из системы");
    } catch (error) {
      console.error("Logout error:", error);
      message.error("Ошибка при логауте");
    } finally {
      localStorage.removeItem("token");
      setToken(null);
      delete axios.defaults.headers.common["Authorization"];
    }
  };

  if (!token) {
    return <LoginPage serverUrl={serverUrl} onLogin={handleLogin} />;
  }

  // Определяем роль пользователя из токена
  let role = "user";
  try {
    const decoded = jwtDecode<JwtPayload>(token);
    role = decoded.role;
  } catch (error) {
    console.error("Ошибка декодирования токена", error);
  }

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Sider width={200}>
        <Menu theme="dark" mode="inline" defaultSelectedKeys={["1"]}>
          {role === "admin" && <Menu.Item key="1">Все устройства</Menu.Item>}
          <Menu.Item key="2">Моё устройство</Menu.Item>
          {/* <Menu.Item key="3" onClick={handleLogout}>
            Выход
          </Menu.Item> */}
        </Menu>
      </Sider>
      <Layout>
        <Header
          style={{
            color: "white",
            fontSize: "20px",
            display: "flex",
            justifyContent: "space-between",
          }}
        >
          <div>MDM Admin Panel</div>
          <Button type="primary" onClick={handleLogout}>
            Выйти
          </Button>
        </Header>
        <Content style={{ padding: "20px" }}>
          {role === "admin" ? (
            <DeviceList serverUrl={serverUrl} />
          ) : (
            <>
              <DeviceCard deviceId="android-test" serverUrl={serverUrl} />
              <HeartbeatLog deviceId="android-test" serverUrl={serverUrl} />
            </>
          )}
        </Content>
        <Footer style={{ textAlign: "center" }}>MDM Admin Panel © 2025</Footer>
      </Layout>
    </Layout>
  );
};
