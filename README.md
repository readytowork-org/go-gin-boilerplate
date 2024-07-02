# Demo API

Demo API template includes all the common packages and setup used for API development in this Company.

## Development

- Copy `.env.example` to `.env` and update according to requirement.
- Create `serviceAccountKey.json` file for firebase admin sdk.
- To run `docker-compose up` (with default configuration will run at `5000` and adminer runs at `5001`)

## Swagger docs config
**Note:** client will generate schema from swagger docs so, please follow these configurations

- ### Struct

  - use **validate:"required"** tag to insure generated fields/keys are not null/optional.
    ```go
    type User struct {
        Name string  `json:"name" validate:"required"`
    }
    ```
    will be generated as
    ```ts
    interface User {
        name: string
    }
    ```
    else
    ```ts
    interface User {
        name?: string
    }
    ```


- ### Function Comments

  - Using `// @Tags` for User and Admin Access endpoints
    - For easy navigation and understanding, utilize distinct @Tags values for endpoints accessible by users and admins. For instance, 
    - use: 
      - ``// @Tags UserApi`` for endpoints under /users accessible to regular users.
      - ``// @Tags UserManagementApi`` for endpoints under /users restricted to administrative access.


## Run CLI 🖥

- Run `docker-compose exec web sh`
- After running type `./__debug_bin cli` you will start cli application.
- Choose the commands to run afterwards.
- To run `docker-compose up` ( with default configuration will run at 5000 and adminer runs at 5001)
- To run with setting up pre-commit hook `make start` ( with default configuration will run at 5000 and adminer runs at 5001`)

#### Migration Commands 🛳

| Command             | Desc                                                 |
|---------------------| ---------------------------------------------------- |
| `make install`      | installs goalngci-lint and change the hooks config   |
| `make start`        | setup pre-commit hook and runs the project           |
| `make run`          | runs the project                                     |
| `make migrate-up`   | runs migration up command                            |
| `make migrate-down` | runs migration down command                          |
| `make force`        | Set particular version but don't run migration       |
| `make goto`         | Migrate to particular version                        |
| `make drop`         | Drop everything inside database                      |
| `make create`       | Create new migration file(up & down)                 |
| `make crud`         | Create crud template                                 |
| `make swag`         | Run this command to generate swag docs               |

### Implemented Feature

- Dependency Injection (go-fx)
- Routing (gin web framework)
- Environment Files
- Logging (file saving on production) zap
- Middlewares (cors)
- Rate Limiting Middleware
- Middleware transaction
- Database Setup (mysql)
- Models Setup and Automigrate (gorm)
- Repositories
- Implementing Basic CRUD Operation
- CRUD Scaffold Generator
- Migration Runner Implementation
- Live code refresh
- Firebase basic setup.
- GCP Cloud Storage bucket setup.
- Twilio Basic setup.
- Gmail Api Setup.

**For Debugging 🐞** Debugger runs at `5002`. Vs code configuration is at `.vscode/launch.json` which will attach debugger to remote application.

## For auto generate of CRUD(Create, ReaD, Update & Delete) api following informations are needed and will be asked in terminal:

- resource-name: name of CRUD in upper camelCase. examples:Food,Puppy,ProductCategory etc.

- resource-table-name: name of CRUD in lower snake case. examples:food,puppy,product_category etc.

- plural-resource-table-name: plural name for the table going to be created. example: foods, puppies, product_categories.

- plural-resource-name: plural name of CRUD in Upper camelCase. examples:Foods,Puppies,ProductCategories etc.
