import React, { useState } from "react";
import { Form, Input, Button, message } from "antd";
import axios from "axios";

interface LoginFormValues {
  username: string;
  password: string;
}

interface LoginPageProps {
  serverUrl: string;
  onLogin: (token: string) => void;
}

const LoginPage: React.FC<LoginPageProps> = ({ serverUrl, onLogin }) => {
  const [loading, setLoading] = useState(false);

  const onFinish = async (values: LoginFormValues) => {
    setLoading(true);
    try {
      const response = await axios.post(`${serverUrl}/login`, values, {
        headers: { "Content-Type": "application/json" },
      });
      const token = response.data.token;
      message.success("Вход выполнен успешно");
      onLogin(token);
    } catch (error: unknown) {
      console.error("Login error:", error);
      message.error("Ошибка входа");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 400, margin: "100px auto" }}>
      <h2>Вход в систему</h2>
      <Form name="login" onFinish={onFinish}>
        <Form.Item
          name="username"
          rules={[{ required: true, message: "Введите имя пользователя" }]}
        >
          <Input placeholder="Имя пользователя" />
        </Form.Item>
        <Form.Item
          name="password"
          rules={[{ required: true, message: "Введите пароль" }]}
        >
          <Input.Password placeholder="Пароль" />
        </Form.Item>
        <Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            loading={loading}
            style={{ width: "100%" }}
          >
            Войти
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
};

export default LoginPage;
