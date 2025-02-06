CREATE TABLE device (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    device_id TEXT NOT NULL,
    camera_enabled BOOLEAN NOT NULL DEFAULT false,
    last_heartbeat TIMESTAMP ,
    created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP
);
