# phlog

phlog on fly

## Running locally

```console
$ go run ./cmd/server
2021/07/08 19:50:55 listening on http://0.0.0.0:8080
```

## Running with postgres

RUN pg in docker:

`docker run -p 5432:5432 -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=phlog -d postgres:13.2`

copy & export .env file

`cp .env.example .env && export $(grep -v '^#' .env | xargs)`

```console
$ go run ./cmd/server
2021/07/08 19:50:55 listening on http://0.0.0.0:8080
```

## Building locally

```console
$ docker build -t phlog .
```

## Deployment

```console
$ fly deploy
```

## Deployment with durable storage

```console
$ fly postgres create
...
$ fly postgres attach --postgres-app {PG_APP} -a {APP}
...
$ fly secrets set -a {APP} AWS_S3_BUCKET="" AWS_S3_KEY="" AWS_S3_SECRET="" AWS_S3_ENDPOINT="" AWS_S3_REGION=""
```
