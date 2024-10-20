DOCKER_IMAGE_APP_NAME=flight_app

build_app:
	docker build -t $(DOCKER_IMAGE_APP_NAME) -f docker/flight/Dockerfile .

deploy_app:
	docker run --name flight_app --rm --network host -p 8080:8080 --env-file docker/config/.env $(DOCKER_IMAGE_APP_NAME) --restart=always

deploy_db:
	docker run --name postgres --rm -e POSTGRES_USER=flight_admin -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_DB=flight_app -p 5432:5432 -d postgres:14.4-alpine 

migrate_up:
	# https://github.com/golang-migrate/migrate
	docker run --rm -v "$(shell pwd)/migrations:/migrations" --network host migrate/migrate -path=/migrations/ -database "postgres://flight_admin@:5432/flight_app?sslmode=disable" up

up:
	docker-compose up -d
