CREATE TABLE IF NOT EXISTS citydrive.cars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brand TEXT NOT NULL,
    model TEXT NOT NULL,
    year_of_manufacture INTEGER NOT NULL,
    fuel_type TEXT NOT NULL CHECK (fuel_type IN ('diesel', '92', '95', '98')),
    license_plate TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);