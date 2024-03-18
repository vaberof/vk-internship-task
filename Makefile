PROJECT_NAME=film-library
APP_LOCAL_NAME=web-backend

DOCKER_LOCAL_IMAGE_NAME=$(PROJECT_NAME)/$(APP_LOCAL_NAME)

WORK_DIR_LINUX=./cmd/filmlibrary
CONFIG_DIR_LINUX=./cmd/filmlibrary/config

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=admin
POSTGRES_DATABASE=film_library

DB_URL="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable"

MIGRATIONS_PATH=migrations/

docker.run: docker.build
	docker compose -f cmd/filmlibrary/docker-compose.yaml up -d

docker.build: build.linux
	docker build -t $(DOCKER_LOCAL_IMAGE_NAME) -f cmd/filmlibrary/Dockerfile .

run.linux: build.linux
	go run $(WORK_DIR_LINUX)/*.go \
		-config.files $(CONFIG_DIR_LINUX)/application.yaml \
		-env.vars.file $(CONFIG_DIR_LINUX)/application.env \

build.linux: build.linux.clean
	mkdir -p $(WORK_DIR_LINUX)/build
	go build -o $(WORK_DIR_LINUX)/build/main $(WORK_DIR_LINUX)/*.go
	cp -R $(CONFIG_DIR_LINUX)/* $(WORK_DIR_LINUX)/build

build.linux.local: build.linux.clean
	mkdir -p $(WORK_DIR_LINUX)/build
	go build -o $(WORK_DIR_LINUX)/build/main $(WORK_DIR_LINUX)/*.go
	cp -R $(CONFIG_DIR_LINUX)/* $(WORK_DIR_LINUX)/build
	@echo "build.local: OK"

build.linux.clean:
	rm -rf $(WORK_DIR_LINUX)/build

migrate.up:
	migrate -path $(MIGRATIONS_PATH) -database $(DB_URL) up

migrate.down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

swagger.gen:
	swag init --parseDependency --parseInternal -g ./cmd/filmlibrary/main.go -o ./cmd/filmlibrary/docs

tests.run:
	go test ./internal/domain/tests/...

tests.run.verbose:
	go test -v ./internal/domain/tests/...

mock.gen: mock.actor_storage.gen mock.film_storage.gen mock.user_finder.gen mock.user_storage.gen

mock.actor_storage.gen:
	mockgen -source=internal/domain/actor_storage.go \
	-destination=internal/domain/mocks/mock_actor_storage.go

mock.film_storage.gen:
	mockgen -source=internal/domain/film_storage.go \
	-destination=internal/domain/mocks/mock_film_storage.go

mock.user_finder.gen:
	mockgen -source=internal/service/auth/user_finder.go \
	-destination=internal/service/auth/mocks/mock_user_finder.go

mock.user_storage.gen:
	mockgen -source=internal/service/user/user_storage.go \
	-destination=internal/service/user/mocks/mock_user_storage.go