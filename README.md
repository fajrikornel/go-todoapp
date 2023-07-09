# To-Do App in Go

Hopefully a CRUD app that takes vitamins.

### Initialize config:

```
cp config/application.yml.sample config/application.yml
```

If you are not using the docker-compose environment, modify the values accordingly.

### Execute DB migrations

Activate the DB and test DB
```
docker-compose up -d db db-test
```

DB migrations
```
go build -o build/migrate cmd/migrate/migrate.go
build/migrate
```

Test environment DB migrations
```
go build -o build/test_migrate cmd/testmigrate/test_migrate.go
build/test_migrate
```

### Run tests

If you are using the docker-compose environment, activate test DB:
```
docker-compose up -d db-test
```

Run the test
```
go test -p 1 -v ./test/...
```

### Start the server


To run it with the docker-compose environment:
```
docker-compose up -d app db
```

To run it with the docker-compose environment and with log monitoring:
```
docker-compose up -d app db elasticsearch kibana filebeat
```

To run it without the docker-compose environment:
```
go build -o build/server cmd/server/server.go
build/server
```

### Shutting down the docker-compose environment

```
docker-compose down
```

### Author notes

In this project, I'm trying to create a service that is developed in a healthy way. Some of the traits I'm trying to achieve:
- Development ease: tools to help local development; migration commands, dockerfiles, docker compose
- Configurability: configure app via a YAML configuration
- Well-documented: Swagger file, README instructions
- Well-logged: logs can be correlated with related logs (via correlation ID)
- Operability: logs can be monitored and queried for operability and debugging
- Clean code: readability, separation of concerns, etc (though this is a bit more subjective/can always be improved)
- Follows convention: in this case, following Go community convention (project layout, naming conventions, etc)
- Equipped with good unit tests with good testing structure/framework

Some of the things in my mind that are not yet implemented:
- Configuration via environment variables
- Instrumentation (metrics, dashboards)
- Deployment tools: CI/CD pipeline
