#!/bin/bash
docker compose up db -d
migrate -database "postgresql://root:example@127.0.0.1:5432/testing?sslmode=disable" -source file://db/migrations up
go clean -testcache
go test -coverprofile cover.out ./pkg/short
go tool cover -html=cover.out