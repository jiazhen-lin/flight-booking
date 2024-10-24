CREATE TABLE IF NOT EXISTS booking_seats(
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    booking_id BIGINT NOT NULL,
    flight_id BIGINT NOT NULL,
    seats INTEGER NOT NULL
);

CREATE INDEX idx_flight_id ON booking_seats (flight_id);
CREATE INDEX idx_booking_id ON booking_seats (booking_id);
