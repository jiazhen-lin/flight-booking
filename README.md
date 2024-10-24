# flight-booking

This is a simple flight booking system, it includes:
- flight search
- flight booking

### Local Deployment

Requirements: docker

1. `make up`
2. List available flights: `curl -i http://localhost:8080/v1/flights/search?departure_airport_id=1&arrival_airport_id=2&departure_date=2024-10-20&limit=10`
3. Book a flight: `curl -i -X POST http://localhost:8080/v1/flights/book -d '{"flights": [{"flight_id": 2, "seats": 1}]}'`

