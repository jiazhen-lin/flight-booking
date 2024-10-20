INSERT INTO flights (created_at, updated_at, number, departure_airport_id, arrival_airport_id, departure_time, duration_seconds, status, total_seats, overbooked_seats, available_seats, price) VALUES
    (now(), now(), 'AA100', 1, 2, '2024-10-20 10:00:00', 3600, 1, 100, 5, 105, 100.99),
    (now(), now(), 'AA101', 2, 1, '2024-10-20 13:00:00', 3600, 1, 100, 5, 105, 200.00),
    (now(), now(), 'AA102', 1, 2, '2024-10-20 16:00:00', 3600, 1, 100, 10, 110, 150.99),
    (now(), now(), 'AA103', 1, 3, '2024-10-20 19:00:00', 3600, 1, 100, 0, 100, 160.00),
    (now(), now(), 'AA104', 1, 2, '2024-10-20 22:00:00', 3600, 1, 100, 0, 1, 170.00)
;