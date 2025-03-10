# Development

You should stay in the ZITADEL root directory to execute the statements in the following chapters.

## Prerequisite

- Buildkit compatible docker installation
- [go](https://go.dev/doc/install)
- [goreleaser](https://goreleaser.com/install/)

Minimum resources:

- CPU's: 2
- Memory: 4Gb

### Installing goreleaser

If you get the error `Failed to fetch https://repo.goreleaser.com/apt/Packages  Certificate verification failed: The certificate is NOT trusted. The certificate chain uses expired certificate.` while installing goreleaser with `apt`, then ensure that ca-certificates are installed:

```sh
sudo apt install ca-certificates
```

### env variables

You can use the default vars provided in [this .env-file](../build/local/local.env) or create your own and update the paths in the [docker compose file](../build/local/docker-compose-local.yml).

## Local Build

Simply run goreleaser to build locally. This will generate all the required files, such as angular and grpc automatically.

```sh
goreleaser build --snapshot --rm-dist --single-target
```

## Production Build & Release

Simply use goreleaser:

```sh
goreleaser release
```

## Run

### Start storage

Use this if you only want to start the storage services needed by ZITADEL. These services are started in background (detached).

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 \
&& docker compose -f ./build/local/docker-compose-local.yml --profile storage up -d
```

**On apple silicon:**
Restart the command (second terminal `docker restart zitadel_<SERVICE_NAME>_1`) if `db` logs `qemu: uncaught target signal 11 (Segmentation fault) - core dumped` or no logs are written from `db-migrations`.

### Initialize the console

Used to set the client id of the console. This step is for local development. If you don't work with a local backend you have to set the client id manually.

You must [initialise the data](###-Initialise-data)) first.

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 \
&& docker compose -f ./build/local/docker-compose-local.yml --profile console-stub up --exit-code-from client-id
```

The command exists as soon as the client id is set.

### Start the Console

The console service is configured for hot reloading. You can also use docker compose for local development.

If you don't use the backend from local you have to configure [the environment.json](../build/local/environment.json) manually.

If you use the local backend ensure that you have [set the correct client id](###-Initialise-frontend).

#### Run console in docker compose

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-local.yml --profile frontend up
```

### Run backend

Used if you want to run the backend locally. It's recommended to [initialise the data](###-Initialise-data) first.

#### Run backend in docker compose

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-local.yml --profile storage --profile backend up
```

#### Run backend locally

##### Export environment variables

```bash
# exports all default env variables
while read line; do
    if [[ $line != #* ]] && [[ ! -z $line ]]; then
        export $line
    fi
done < build/local/local.env
```

##### Start command for backend

```bash
# starts zitadel with default config files
go run cmd/zitadel/main.go -console=false -localDevMode=true -config-files=cmd/zitadel/startup.yaml -config-files=cmd/zitadel/system-defaults.yaml -config-files=cmd/zitadel/authz.yaml start
```

If you want to run your backend locally and the frontend by docker compose you have to replace the following variables:

[docker compose yaml](../build/local/docker-compose-local.yml):

```yaml
service:
  client-id:
    environment:
      - HOST=backend-run
  grpc-web-gateway:
    environment:
      - BKD_HOST=backend-run
```

with

```yaml
service:
  client-id:
    environment:
      - HOST=host.docker.internal
  grpc-web-gateway:
    environment:
      - BKD_HOST=host.docker.internal
```

##### Setup ZITADEL

The following command starts the backend of ZITADEL with the default config files:

```bash
go run cmd/zitadel/main.go -setup-files=cmd/zitadel/setup.yaml -setup-files=cmd/zitadel/system-defaults.yaml -setup-files=cmd/zitadel/authz.yaml setup
```

## Initial login credentials

**username**: `zitadel-admin@caos-ag.zitadel.ch`

**password**: `Password1!`
