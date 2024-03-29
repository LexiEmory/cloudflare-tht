# Cloudflare Take Home Test

A very simple link shortener endpoint.

# Dependencies

- [Docker](https://docs.docker.com/)
- [Go Migrate Cli](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

# Getting Started

Create a local env file from the example

```shell
$ cp .env.example .env
```

Start up containers

```shell
$ docker-compose up -d --build --force-recreate
```

Migrate to database

```shell
$ migrate -database "postgresql://root:example@127.0.0.1:5432/postgres?sslmode=disable" -source file://db/migrations up
```

All endpoints are documented in the `docs.yaml` file found in the `api` directory.

# Testing

Start by ensuring your database is running.

```shell
$ docker compose up db -d
```

Then run

```shell
$ ./test.sh
```

This script will ensure the docker database is migrated and open a coverage web page.