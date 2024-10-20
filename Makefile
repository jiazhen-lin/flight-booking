DOCKER_IMAGE_APP_NAME=flight_app

build_app:
	docker build -t $(DOCKER_IMAGE_APP_NAME) -f docker/flight/Dockerfile .

run_app:
	docker run --name flight_app --rm -p 8080:8080 --env-file docker/config/.env $(DOCKER_IMAGE_APP_NAME) --restart=always

