import React, { useState, useEffect } from "react";
import { Table, Dropdown, Menu, Button, message } from "antd";
import axios from "axios";

export interface Device {
  id: string;
  device_id: string;
  camera_enabled: boolean;
  microphone_enabled: boolean;
  bluetooth_enabled: boolean;
  os_version: string;
  battery_level: number;
  last_heartbeat: string;
  created_at: string;
  updated_at: string;
}

interface DeviceListProps {
  serverUrl: string;
}

const DeviceList: React.FC<DeviceListProps> = ({ serverUrl }) => {
  const [devices, setDevices] = useState<Device[]>([]);

  const fetchDevices = async () => {
    try {
      const response = await axios.get<Device[]>(`${serverUrl}/devices`);
      setDevices(response.data);
    } catch (error: unknown) {
      console.error("Error fetching devices: ", error);
      message.error("Ошибка при получении списка устройств");
    }
  };

  useEffect(() => {
    fetchDevices();
  }, [serverUrl]);

  // Функция для отправки команды для конкретного устройства
  const sendCommandForDevice = async (
    deviceId: string,
    endpoint: string,
    payload: object
  ) => {
    try {
      await axios.post(`${serverUrl}/devices/${deviceId}/${endpoint}`, payload);
      message.success("Команда успешно отправлена");
      fetchDevices(); // Обновляем список после выполнения команды
    } catch (error) {
      console.error(`Error sending command to ${endpoint}:`, error);
      message.error("Ошибка отправки команды");
    }
  };

  const columns = [
    { title: "Device ID", dataIndex: "device_id", key: "device_id" },
    {
      title: "Камера",
      dataIndex: "camera_enabled",
      key: "camera_enabled",
      render: (val: boolean) => (val ? "Включена" : "Выключена"),
    },
    {
      title: "Микрофон",
      dataIndex: "microphone_enabled",
      key: "microphone_enabled",
      render: (val: boolean) => (val ? "Включен" : "Выключен"),
    },
    {
      title: "Bluetooth",
      dataIndex: "bluetooth_enabled",
      key: "bluetooth_enabled",
      render: (val: boolean) => (val ? "Включен" : "Выключен"),
    },
    { title: "Версия ОС", dataIndex: "os_version", key: "os_version" },
    {
      title: "Заряд",
      dataIndex: "battery_level",
      key: "battery_level",
      render: (val: number) => `${val}%`,
    },
    {
      title: "Последний heartbeat",
      dataIndex: "last_heartbeat",
      key: "last_heartbeat",
      render: (val: string) => new Date(val).toLocaleString(),
    },
    {
      title: "Действия",
      key: "actions",
      render: (_: unknown, record: Device) => (
        <Dropdown
          overlay={
            <Menu>
              <Menu.Item
                onClick={() =>
                  sendCommandForDevice(record.device_id, "camera", {
                    enabled: !record.camera_enabled,
                  })
                }
              >
                Переключить камеру
              </Menu.Item>
              <Menu.Item
                onClick={() =>
                  sendCommandForDevice(record.device_id, "microphone", {
                    enabled: !record.microphone_enabled,
                  })
                }
              >
                Переключить микрофон
              </Menu.Item>
              <Menu.Item
                onClick={() =>
                  sendCommandForDevice(record.device_id, "bluetooth", {
                    enabled: !record.bluetooth_enabled,
                  })
                }
              >
                Переключить Bluetooth
              </Menu.Item>
            </Menu>
          }
        >
          <Button>Действия</Button>
        </Dropdown>
      ),
    },
  ];

  return (
    <div>
      <h2>Все устройства</h2>
      <Table dataSource={devices} columns={columns} rowKey="id" />
    </div>
  );
};

export default DeviceList;
