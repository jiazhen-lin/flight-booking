CREATE INDEX idx_airport_departure_time ON flights (
    departure_airport_id, 
    arrival_airport_id, 
    departure_time
);