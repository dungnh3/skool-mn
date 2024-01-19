IMAGE_NAME_SERVICE = skool_mn
VERSION_SERVICE = 1.0

HOST = 127.0.0.1
PORT = 3306
DATABASE = skool_mn
USER = root
PASSWORD = secret

MOCK_TESTS = $(shell find ./test -type f -name "mock_test.go")
E2E_TESTS = $(shell find ./test -type f -name "e2e_test.go")

migrate-up:
	migrate -source "file://migrations" -database "mysql://$(USER):$(PASSWORD)@tcp($(HOST):$(PORT))/$(DATABASE)" up
.PHONY: migration-up

build:
	go build -o bin/runtime main.go
.PHONY: build

build-docker:
	docker build -t $(IMAGE_NAME_SERVICE):$(VERSION_SERVICE) --force-rm -f Dockerfile .
.PHONY: build-docker

start:
	docker run -it --name $(IMAGE_NAME_SERVICE) -p 9090:9090 $(IMAGE_NAME_SERVICE):$(VERSION_SERVICE)
.PHONY: start

start-docker-compose:
	docker-compose up -d
.PHONY: start-docker-compose

mock:
	cd internal && rm -r mocks && mockery --all --keeptree --case underscore && cd ..
.PHONY: mock

mock-test:
	go test --race -v -count=1 $(MOCK_TESTS)
.PHONY: mock-test

test-e2e:
	go test --race -v -count=1 $(E2E_TESTS)
.PHONY: test-e2e

test-integration: init-db-test test-e2e stop-test-db
.PHONY: test-integration

create-db:
	docker run --rm -d --name database -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret mysql:8.0.31 || true
.PHONY: create-db

init-db: create-db
	docker exec -it database mysql -u root -psecret -e "CREATE DATABASE skool_mn;" || true
.PHONY: init-db

init-db-test: create-db waiting-start-container
	docker exec -it database mysql -u root -psecret -e "CREATE DATABASE skool_mn_test;" || true
	migrate -source "file://migrations" -database "mysql://$(USER):$(PASSWORD)@tcp($(HOST):$(PORT))/skool_mn_test" up
.PHONY: init-db-test

stop-test-db:
	docker exec -it database mysql -u root -psecret -e "DROP DATABASE mfv_test;"
.PHONY: stop-test-db

waiting-start-container:
	echo "waiting init mysql container"
	sleep 10
.PHONY: waiting-start-container

install-swagger:
	go install -v github.com/swaggo/swag/cmd/swag@v1.16.1
.PHONY: install-swagger

swagger: install-swagger
	rm docs/docs.go || true && swag init
.PHONY: swagger