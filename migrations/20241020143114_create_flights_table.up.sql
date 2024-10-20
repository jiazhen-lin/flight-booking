CREATE TABLE IF NOT EXISTS flights(
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    number VARCHAR(10) NOT NULL,
    departure_airport_id INTEGER NOT NULL,
    arrival_airport_id INTEGER NOT NULL,
    departure_time TIMESTAMP NOT NULL,
    duration_seconds INTEGER NOT NULL,
    status SMALLINT NOT NULL,
    total_seats integer NOT NULL,
    overbooked_seats integer NOT NULL,
    available_seats integer NOT NULL,
    price NUMERIC(10,4) NOT NULL
);
