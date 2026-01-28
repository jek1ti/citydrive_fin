CREATE TABLE IF NOT EXISTS citydrive.car_telemetry_history (
    id BIGSERIAL PRIMARY KEY,
    car_id UUID NOT NULL REFERENCES citydrive.cars(id) ON DELETE CASCADE,
    lat DOUBLE PRECISION NOT NULL,
    lon DOUBLE PRECISION NOT NULL,
    fuel REAL NOT NULL CHECK (fuel >= 0 AND fuel <= 100),
    speed REAL NOT NULL CHECK (speed >= 0),
    engine_on BOOLEAN NOT NULL,
    locked BOOLEAN NOT NULL,
    activated BOOLEAN NOT NULL,
    rpm INTEGER NOT NULL CHECK (rpm >= 0),
    handbrake BOOLEAN NOT NULL,
    odo INTEGER NOT NULL CHECK (odo >= 0),
    timestamp BIGINT NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cth_car_id_timestamp ON citydrive.car_telemetry_history(car_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_cth_timestamp ON citydrive.car_telemetry_history(timestamp);
CREATE INDEX IF NOT EXISTS idx_cth_activated ON citydrive.car_telemetry_history(activated);