# Go project template

[![LICENSE](https://img.shields.io/badge/License-MIT-turquise.svg)](LICENSE)

## Overview

The purpose of this template is to have a template to quickstart new projects. It also serves as an example on how to organize a project using hexagonal architecture (ports/adapters) with dummy implementation.

### Philosophy

It is designed to be a simple starting point, using mostly stdlib not to force any lib/framework on the resulting project.

## How to use

### Starting the server

For the moment, juste run `go run cmd/main.go`. Future QoL updates will include containerization and live reload.

### Tests

The template includes tests. Run the tests with `go test ./...`. You can test coverage with `go test ./... -coverprofile=cover.out >> /dev/null && go tool cover -func=cover.out`

<!--### Mocks

Some tests uses mocks. You can generate them using `go generate ./...`. Generate commands in the project includes [mockgen](https://github.com/uber-go/mock) and [gofilemerge](https://github.com/yoannduc/gofilemerge).-->

## Project structure

### `cmd/main.go`

Project initialisation, dependency injection & server start.

### `internal/domain`

Domain scope elements.

### `internal/dtos`

Data transfer object for inputs & outputs.

### `internal/handlers`

Handlers for input.

### `internal/ports`

Ports (interfaces) defining services (core) & repositories.

### `internal/repositories`

Repositories.

### `internal/services`

Services (core).

### `pkg/logger`

Logger.

### `pkg/mapper`

Generic mapper to transform dto to domain & domain to dto.

### `pkg/middlewares`

Middlewares for HTTP API.
