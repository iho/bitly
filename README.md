# bitly

## Description

Simple URL shortener microservice using JSON REST API.

## Installation

```bash
docker-compose up -d
```

## Usage

```bash
curl -X POST -H "Content-Type: application/json" -d '{"url":"https://www.google.com"}' http://localhost:8000/urls
curl http://localhost:8080/urls/ytKEQIFa/
```

## Run tests

```bash
make gen
go test ./... -v
```

Running tests can be slow since I use dockertest from Ory and it starts a new docker container for each test.
