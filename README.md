# To-Do App in Go

Hopefully a CRUD app that takes vitamins.

### Initialize config:

```
cp config/application.yml.sample config/application.yml
```

If you are not using the docker-compose environment, modify the values accordingly.

The configs set on `config/application.yml` can be overriden with environment variables.
Set the environment variable with `TODOAPP_` as a prefix. For nested structures exampled by the YAML config, replace `.` with `_` (`DB_CONFIG.DB_HOST` becomes `TODOAPP_DB_CONFIG_DB_HOST`)

Log and metrics monitoring config rest in the `filebeat.docker.yml` and `prometheus.yml` in the /config folder respectively.

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
docker-compose up -d --build app db
```

To run it with the docker-compose environment and with log monitoring:
```
docker-compose up -d --build app db elasticsearch kibana filebeat
```

To run it with the docker-compose environment and with metrics monitoring:
```
docker-compose up -d --build app db prometheus grafana
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

### Usage

Using the app:
```
curl localhost:8080/v1/projects/1
```
For API details, see the Swagger file in the api/ folder of this repo.

Provided that the monitoring stack is turned on (see above instructions),

- Log monitoring:
```
1. Visit Kibana on localhost:5601 using a browser
2. Go to side pane > Analytics > Discover
3. Observe or query the logs
```

- Metrics monitoring:
```
1. Visit Grafana UI on localhost:3000
2. For first time login, use `admin` user with the password as `admin` and change the password.
3. Add a new Grafana data source as Prometheus type with the URL as http://prometheus:9090
4. You can then add dashboards according to the metrics available in the app
```


### Author notes

In this project, I'm trying to create a service that is developed in a healthy way. Some of the traits I'm trying to achieve:
- Development ease: tools to help local development; migration commands, dockerfiles, docker compose
- Configurability: configure app via a YAML configuration and override it via environment variables
- Instrumented with metrics and dashboards to monitor application-specific health
- Well-documented: Swagger file, README instructions
- Well-logged: logs can be correlated with related logs (via correlation ID)
- Operability: logs can be monitored and queried for operability and debugging
- Clean code: readability, separation of concerns, etc (though this is a bit more subjective/can always be improved)
- Follows convention: in this case, following Go community convention (project layout, naming conventions, etc)
- Equipped with good unit tests with good testing structure/framework

Some of the things in my mind that are not yet implemented:
- Deployment tools: CI/CD pipeline
