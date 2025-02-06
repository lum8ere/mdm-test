CREATE TABLE device (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    device_id TEXT NOT NULL,
    camera_enabled BOOLEAN NOT NULL DEFAULT false,
    microphone_enabled BOOLEAN NOT NULL DEFAULT false,
    bluetooth_enabled BOOLEAN NOT NULL DEFAULT false,
    os_version TEXT,
    battery_level INT, 
    last_heartbeat TIMESTAMP ,
    created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP
);

CREATE TABLE users (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL, -- Например, 'admin' или 'user'
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
