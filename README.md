# Go Gin Skeleton

> **Go Gin Skeleton template includes all the common packages and setup used for API development using [gin](https://gin-gonic.com).**

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
  - migrate: for DB migrations
  - gentool: to generate dao objects from database
  - swag: to generate swagger docs
  - gin: hot-reload

**For Debugging 🐞** Debugger runs at `5002`. Vs code configuration is at `.vscode/launch.json` which will attach debugger to remote application.

## Development

- Copy `.env.example` to `.env` and update according to requirement.
- #### Running using Docker
  - run `docker-compose up` (with default configuration will run at `5000` and adminer runs at `5001`)
- #### Running using Gin Watch
  - run `make dev`

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
- To run with setting up pre-commit hook `make start` ( with default configuration will run at 5000 and adminer runs at 5001`)

## Implements Google Cloud Proxy by default

This reduces hassle for developer to update IP in Cloud SQL during IP Change.

Implemented through docker image as follows

- Cloud SQL -> Google Proxy Docker Image -> Web App

#### Points to remember for smooth working

- `ServiceAccountKey.json` requires Cloud SQL Read and Write Permission
- `DB_HOST_NAME` value in `.env` is required
- `DB_HOST=cloud-sql-proxy` instead of `IPV4` or `DB_HOST_NAME` for development environment
- `DB_PORT` will be `3306` by default

#### Migration Commands 🛳

| Command             | Desc                                               |
| ------------------- | -------------------------------------------------- |
| `make install`      | installs goalngci-lint and change the hooks config |
| `make start`        | setup pre-commit hook and runs the project         |
| `make run`          | runs the project                                   |
| `make migrate-up`   | runs migration up command                          |
| `make migrate-down` | runs migration down command                        |
| `make force`        | Set particular version but don't run migration     |
| `make goto`         | Migrate to particular version                      |
| `make drop`         | Drop everything inside database                    |
| `make create`       | Create new migration file(up & down)               |
| `make crud`         | Create crud template                               |
| `make swag`         | Run this command to generate swag docs             |

## For auto generate of CRUD(Create, ReaD, Update & Delete) api following informations are needed and will be asked in terminal:

- resource-name: name of CRUD in upper camelCase. examples:Food,Puppy,ProductCategory etc.

- resource-table-name: name of CRUD in lower snake case. examples:food,puppy,product_category etc.

- plural-resource-table-name: plural name for the table going to be created. example: foods, puppies, product_categories.

- plural-resource-name: plural name of CRUD in Upper camelCase. examples:Foods,Puppies,ProductCategories etc.
