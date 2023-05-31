# go-todoapp
Go To-do App for practicing Go

### Initialize config:

```
cp config/application.yml.sample config/application.yml
```

### Activate local database:

```
docker-compose up -d
```

### Execute DB migrations

```
go build -o build/migrate cmd/migrate/migrate.go
build/migrate
```

### Start the server

```
go build -o build/server cmd/server/server.go
build/server
```
