import React, { useState, useEffect } from 'react';
import { List } from 'antd';
import axios from 'axios';

interface HeartbeatLogProps {
  deviceId: string;
  serverUrl: string;
}

const HeartbeatLog: React.FC<HeartbeatLogProps> = ({ deviceId, serverUrl }) => {
  const [logs, setLogs] = useState<string[]>([]);

  // Функция для получения heartbeat через REST API
  const fetchHeartbeat = async () => {
    try {
      // Предполагаем, что endpoint heartbeat возвращает данные устройства
      const response = await axios.post(`${serverUrl}/devices/${deviceId}/heartbeat`);
      const logEntry = `Heartbeat: ${new Date(response.data.last_heartbeat).toLocaleString()}`;
      setLogs((prev) => [logEntry, ...prev]);
    } catch (error) {
      console.error('Error fetching heartbeat:', error);
    }
  };

  // Запускаем polling heartbeat каждые 10 секунд
  useEffect(() => {
    fetchHeartbeat();
    const interval = setInterval(fetchHeartbeat, 10000);
    return () => clearInterval(interval);
  }, [serverUrl, deviceId]);

  return (
    <div style={{ marginTop: '20px' }}>
      <h3>Лог Heartbeat</h3>
      <List
        bordered
        dataSource={logs}
        renderItem={(item) => <List.Item>{item}</List.Item>}
      />
    </div>
  );
};

export default HeartbeatLog;
