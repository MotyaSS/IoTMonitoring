CREATE TABLE IF NOT EXISTS devices (
    device_id TEXT PRIMARY KEY,
    wind_speed_high_bound DOUBLE PRECISION NOT NULL,
    wind_speed_low_bound DOUBLE PRECISION NOT NULL,
    temperature_high_bound DOUBLE PRECISION NOT NULL,
    temperature_low_bound DOUBLE PRECISION NOT NULL,
    pressure_high_bound DOUBLE PRECISION NOT NULL,
    pressure_low_bound DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_devices_created_at ON devices (created_at);

INSERT INTO devices (
    device_id,
    wind_speed_high_bound,
    wind_speed_low_bound,
    temperature_high_bound,
    temperature_low_bound,
    pressure_high_bound,
    pressure_low_bound
) VALUES
    ('1', 30.0, 0.0, 70.0, -30.0, 1100.0, 900.0),
    ('2', 25.0, 0.0, 60.0, -20.0, 1080.0, 920.0)
ON CONFLICT (device_id) DO NOTHING;

