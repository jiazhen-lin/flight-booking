CREATE TABLE IF NOT EXISTS bookings(
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    flight_id UUID NOT NULL,
    seats INTEGER NOT NULL,
    status SMALLINT NOT NULL
);
