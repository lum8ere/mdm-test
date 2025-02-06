import React, { useState, useEffect } from 'react';
import { Card, Button, Dropdown, Menu, message } from 'antd';
import axios from 'axios';

interface Device {
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

interface DeviceCardProps {
  deviceId: string;
  serverUrl: string;
}

const DeviceCard: React.FC<DeviceCardProps> = ({ deviceId, serverUrl }) => {
  const [device, setDevice] = useState<Device | null>(null);

  // Функция для получения статуса устройства через REST API
  const fetchStatus = async () => {
    try {
      const response = await axios.get<Device>(`${serverUrl}/devices/${deviceId}/status`);
      setDevice(response.data);
    } catch (error) {
      console.error('Error fetching device status:', error);
    }
  };

  // Запускаем polling каждые 5 секунд
  useEffect(() => {
    fetchStatus();
    const interval = setInterval(fetchStatus, 5000);
    return () => clearInterval(interval);
  }, [serverUrl, deviceId]);

  // Функция для отправки команды (например, переключение состояния)
  const sendCommand = async (endpoint: string, payload: object) => {
    try {
      await axios.post(`${serverUrl}/devices/${deviceId}/${endpoint}`, payload);
      message.success('Команда успешно отправлена');
      fetchStatus(); // Обновляем статус после отправки команды
    } catch (error) {
      console.error(`Error sending command to ${endpoint}:`, error);
      message.error('Ошибка отправки команды');
    }
  };

  // Меню для управления устройством
  const menu = (
    <Menu>
      <Menu.Item onClick={() => sendCommand('camera', { enabled: !device?.camera_enabled })}>
        Переключить камеру
      </Menu.Item>
      <Menu.Item onClick={() => sendCommand('microphone', { enabled: !device?.microphone_enabled })}>
        Переключить микрофон
      </Menu.Item>
      <Menu.Item onClick={() => sendCommand('bluetooth', { enabled: !device?.bluetooth_enabled })}>
        Переключить Bluetooth
      </Menu.Item>
    </Menu>
  );

  return (
    <Card
      title={`Устройство: ${deviceId}`}
      extra={
        <Dropdown overlay={menu} trigger={['click']}>
          <Button>Управление</Button>
        </Dropdown>
      }
      style={{ marginBottom: '20px' }}
    >
      {device ? (
        <div>
          <p>
            <strong>Камера:</strong> {device.camera_enabled ? 'Включена' : 'Выключена'}
          </p>
          <p>
            <strong>Микрофон:</strong> {device.microphone_enabled ? 'Включен' : 'Выключен'}
          </p>
          <p>
            <strong>Bluetooth:</strong> {device.bluetooth_enabled ? 'Включен' : 'Выключен'}
          </p>
          <p>
            <strong>Версия ОС:</strong> {device.os_version || 'N/A'}
          </p>
          <p>
            <strong>Уровень заряда:</strong> {device.battery_level}%
          </p>
          <p>
            <strong>Последний heartbeat:</strong>{' '}
            {new Date(device.last_heartbeat).toLocaleString()}
          </p>
        </div>
      ) : (
        <p>Загрузка...</p>
      )}
    </Card>
  );
};

export default DeviceCard;
