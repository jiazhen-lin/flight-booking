INSERT INTO flights (created_at, updated_at, number, departure_airport_id, arrival_airport_id, departure_time, duration_seconds, status, total_seats, overbooked_seats, available_seats, price) VALUES
    (now(), now(), 'AA101', 1, 3, '2024-10-20 02:00:00', 21600, 1, 100, 0, 1, 50.99),
    (now(), now(), 'AA100', 1, 2, '2024-10-20 02:00:00', 36000, 1, 100, 1, 101, 130.99),
    (now(), now(), 'AA103', 3, 2, '2024-10-20 11:00:00', 21600, 1, 100, 0, 100, 49.99),
    (now(), now(), 'AA102', 1, 2, '2024-10-20 16:00:00', 36000, 1, 100, 1, 101, 150.99),
    (now(), now(), 'AA103', 2, 1, '2024-10-20 21:00:00', 18000, 1, 100, 0, 100, 160.00),
    (now(), now(), 'AA104', 1, 2, '2024-10-21 00:00:00', 36000, 1, 100, 0, 100, 170.00)
;