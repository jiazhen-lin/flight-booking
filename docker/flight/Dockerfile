FROM golang:1.22-alpine
LABEL maintainer="Jan <jan.jiazhen@gmail.com>"

RUN mkdir -p /app
WORKDIR /app
COPY . .
RUN go build -o flight_app ./cmd/main.go

EXPOSE 8080
ENTRYPOINT ["./flight_app"]
