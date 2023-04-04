# gogql - Backend

## This projects is based upon Golang v1.20

## Description
This is an of implementation of Clean Architecture in Go (Golang) projects.

### Project Structure
This project has  4 Domain layer :
 * Models Layer --> models
 * Repository Layer --> store
 * Master Layer --> master
 * Usecase Layer   --> services
 * Delivery Layer --> api

_Dependencies_

- [go](https://golang.org/)
- [postgres](https://www.postgresql.org/)
- [golang-migrate](https://github.com/golang-migrate/migrate/releases)
- [gqlgen](https://gqlgen.com/)
- [docker](https://docs.docker.com/install/linux/docker-ce/ubuntu/)
- [docker-compose](https://docs.docker.com/compose/install/)

_Included dependent binaries_

- [migrate](https://github.com/golang-migrate/migrate)


## Prerequisites

### Setup docker and docker compose
- Docker installation guide - https://docs.docker.com/engine/install/ubuntu/
- Docker Compose installation guide - https://docs.docker.com/compose/install/


### Setup Golang: go version go1.20.2 linux/amd64
- Golang setup guide - https://golang.org/doc/install
- `wget -c https://golang.org/dl/go1.20.2.linux-amd64.tar.gz`
- `sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.20.2.linux-amd64.tar.gz`
- `export PATH=$PATH:/usr/local/go/bin`
- Setup go path: in root directory, go to file .profile and paste the following line `export PATH=$PATH:/usr/local/go/bin`

#### Setup Golang with gobrew: alternative
- `curl -sLk https://git.io/gobrew | sh -`
- add following command in .bashrc file `export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"`

#### Setup golang-migrate
- `curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash`
- `apt-get update`
- `apt-get install -y migrate`


### Database

#### Setup Postgres & Run Server
```bash

make postgres

make createdb

make migrateup

make dbseed

make run

```

#### Create and Migrate DB
```bash
make createdb
make migratedown
make migrateup


```

#### Seed DB
```bash
make dbseed
```

#### Drop DB
```bash
make migratedown
make dropdb
```

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
\q
```

#### Generate GraphQL
```bash
make gqlgen
```
### Server

#### Test server

```bash
make test
```

#### Run server

```bash
make run
```

#### Build server

```bash
make build
```

## Docker

#### Build image

```bash
docker build -t gogql:latest .
```
