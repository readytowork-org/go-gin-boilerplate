# Go Gin Skeleton

> **Go Gin Skeleton template includes all the common packages and setup used for API development
using [gin](https://gin-gonic.com).**

## Configured With

- Dependency Injection: [fx](https://github.com/uber-go/fx)
- Routing: [gin web framework](https://gin-gonic.com)
- Logging: [zap](https://github.com/uber-go/zap)
- Database: ([mysql](https://gorm.io/driver/mysql) / [sqlmock](https://github.com/DATA-DOG/go-sqlmock))
- ORM: [gorm](https://gorm.io/docs)
- API documentation: [gin-swagger](https://github.com/swaggo/gin-swagger)
- Middlewares
  - CORS
  - Rate Limit
  - DB Transaction
- CLI tools
  - [migrate](https://github.com/golang-migrate/migrate): for DB migrations
  - [gentool](https://gorm.io/gen/): to generate dao objects from database
  - [swag](https://github.com/swaggo/swag): to generate swagger docs
  - [gin](https://github.com/codegangsta/gin): hot-reload

**For Debugging 🐞** Debugger runs at `5002`. Vs code configuration is at `.vscode/launch.json` which will attach
debugger to remote application.

## Running

- Copy `.env.example` to `.env` and update according to requirement.
- ### using Docker
  - run `docker-compose up` (with default configuration will run at `5000` and adminer runs at `5001`)
- ### using Gin Watch
  - run `make run`

## Commands 🛳

| Command                  | Desc                                                                                         |
|--------------------------|----------------------------------------------------------------------------------------------|
| `make install`           | installs goalngci-lint and change the hooks config                                           |
| `make run`               | runs the project using gin watcher                                                           |
| `make migrate <options>` | runs [migrate](https://github.com/golang-migrate/migrate) command with env configs from .env |
| `make crud`              | Create crud template                                                                         |
| `make swagger`           | Run this command to generate swag docs                                                       |
| `make dao`               | Generates go structs from database                                                           |

[//]: # (TODO :: Need a proper name ⬇️)

## External Services

> Run `go get <package name>` to install.

#### Firebase

- Package name: github.com/readytowork-org/go_firebase_service
- Github: https://github.com/readytowork-org/go_firebase_service

#### GCP

- Package name: github.com/readytowork-org/go_gcp_service
- Github: https://github.com/readytowork-org/go_gcp_service

## Swagger docs config

> Please refer to [SWAGGER.md](https://github.com/readytowork-org/go-gin-skeleton/blob/develop/SWAGGER.md)

## Run CLI 🖥

- Run `docker-compose exec web sh`
- After running type `./__debug_bin cli` you will start cli application.
- Choose the commands to run afterwards.
- To run `docker-compose up` ( with default configuration will run at 5000 and adminer runs at 5001)
- To run with setting up pre-commit hook `make start` ( with default configuration will run at 5000 and adminer runs at
  5001`)

## Implements Google Cloud Proxy by default

This reduces hassle for developer to update IP in Cloud SQL during IP Change.

Implemented through docker image as follows

- Cloud SQL -> Google Proxy Docker Image -> Web App

#### Points to remember for smooth working

- `ServiceAccountKey.json` requires Cloud SQL Read and Write Permission
- `DB_HOST_NAME` value in `.env` is required
- `DB_HOST=cloud-sql-proxy` instead of `IPV4` or `DB_HOST_NAME` for development environment
- `DB_PORT` will be `3306` by default

## For auto generate of CRUD(Create, ReaD, Update & Delete) api following informations are needed and will be asked in terminal:

- resource-name: name of CRUD in upper camelCase. examples:Food,Puppy,ProductCategory etc.

- resource-table-name: name of CRUD in lower snake case. examples:food,puppy,product_category etc.

- plural-resource-table-name: plural name for the table going to be created. example: foods, puppies,
  product_categories.

- plural-resource-name: plural name of CRUD in Upper camelCase. examples:Foods,Puppies,ProductCategories etc.
