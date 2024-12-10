include .env

DB_DSN="${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"

# host and port used are based on docker config
DB_DSN_DOCKER="${DB_USERNAME}:${DB_PASSWORD}@tcp(localhost:33066)/${DB_NAME}"

MIGRATE_LOCAL=migrate -path=database/migration -database ${DB_TYPE}"://"${DB_DSN} -verbose

MIGRATE=docker-compose exec web ${MIGRATE_LOCAL}

GEN_TOOL=gentool -fieldNullable -fieldWithIndexTag -fieldWithTypeTag -fieldSignable -onlyModel -outPath './database/dao' -modelPkgName 'dao'

migrate:
         ifeq (migrate,$(firstword $(MAKECMDGOALS)))
           RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
           # ...and turn them into do-nothing targets
           $(eval $(RUN_ARGS):;@:)
         endif
         ifeq (create,$(firstword $(RUN_ARGS)))
           ARGS := $(wordlist 2,$(words $(RUN_ARGS)),$(RUN_ARGS))
           RUN_ARGS=create -ext sql -dir database/migration $(ARGS)
         endif
migrate:
		@echo "using database: ${DB_NAME}"
		@if [ "$(env)" = "local" ]; then $(MIGRATE_LOCAL) $(RUN_ARGS); else $(MIGRATE) $(RUN_ARGS); fi

dao:
		@command -v gentool >/dev/null 2>&1 || (echo "Installing gentool..." && go install gorm.io/gen/tools/gentool@latest)
		@if [ "$(env)" = "local" ]; then $(GEN_TOOL) -dsn $(DB_DSN); else $(GEN_TOOL) -dsn $(DB_DSN_DOCKER); fi

swagger:
		@command -v swag >/dev/null 2>&1 || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
		swag fmt
		swag init --output ./swagger --parseDependency --parseInternal

crud:
		bash automate/scripts/crud.sh

lint-install:
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2
		git config core.hooksPath hooks

run:
		@command -v gin >/dev/null 2>&1 || (go install github.com/codegangsta/gin@latest);
		gin -a $(SERVER_PORT) -i -p $$(($(SERVER_PORT) + 1)) run .

test-repo: TEST_NAME=$(filter-out $@,$(MAKECMDGOALS))
test-repo:
	go test ./tests/repository_test -v -run $(TEST_NAME)

i-test-controller: TEST_NAME=$(filter-out $@,$(MAKECMDGOALS))
i-test-controller:
	go test ./tests/controllers_i_test -v -run $(TEST_NAME)

test-controller: TEST_NAME=$(filter-out $@,$(MAKECMDGOALS))
test-controller:
	go test ./tests/controllers_test -v -run $(TEST_NAME)

.PHONY: dao migrate create swagger test-repo lint-install

context-upload:
	bash automate/scripts/ci-upload.sh
