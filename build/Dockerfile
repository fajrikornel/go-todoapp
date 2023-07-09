FROM golang:1.20.5-alpine3.18

WORKDIR /app

COPY .. /app
COPY ../config/application.yml.sample /app/config/application.yml

RUN go build -o /app/build/server /app/cmd/server/server.go

CMD ["/bin/sh", "-c", "/app/build/server"]
