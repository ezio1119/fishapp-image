CWD = $(shell pwd)
DC_FILE = docker-compose.yml
SVC = image
DB_SVC = image-db
DB_NAME = image_DB
DB_USER = root
DB_PWD = password
NET = fishapp-net
GRPC_SVC = ImageService
PJT_NAME = $(notdir $(PWD))
# TEST = $(shell docker inspect $(NET) > /dev/null 2>&1; echo " $$?")

createnet:
	docker network create $(NET)

sqldoc: migrate
	docker run --rm --name tbls --net $(NET) -v $(CWD)/db:/work ezio1119/tbls \
	doc -f -t svg mysql://$(DB_USER):$(DB_PWD)@$(DB_SVC):3306/$(DB_NAME) ./

proto:
	docker run --rm --name protoc -v $(CWD)/pb:/pb -v $(CWD)/schema:/proto ezio1119/protoc \
	-I/proto \
	-I/go/src/github.com/envoyproxy/protoc-gen-validate \
	--go_out=plugins=grpc:/pb \
	--validate_out="lang=go:/pb" \
	image.proto

newsql:
	docker run --rm --name newsql -v $(CWD)/db/sql:/sql \
	migrate/migrate:latest create -ext sql -dir ./sql ${a}

test:
	docker-compose -f $(DC_FILE) exec $(SVC) sh -c "go test -v -coverprofile=cover.out ./... && \
	go tool cover -html=cover.out -o ./cover.html" && \
	open ./src/cover.html

cli:
	docker run --rm --name grpc_cli --net $(NET) znly/grpc_cli \
	call $(SVC):50051 $(SVC).$(GRPC_SVC).$(m) "$(q)"

waitdb: updb
	docker run --rm --name dockerize --net $(NET) jwilder/dockerize \
	-timeout 60s \
	-wait tcp://$(DB_SVC):3306

migrate: waitdb
	docker run --rm --name migrate --net $(NET) \
	-v $(CWD)/db/sql:/sql migrate/migrate:latest \
	-path /sql/ -database "mysql://$(DB_USER):$(DB_PWD)@tcp($(DB_SVC):3306)/$(DB_NAME)" $(MIGRATE)

up: migrate
	docker-compose -f $(DC_FILE) up -d $(SVC)

updb:
	docker-compose -f $(DC_FILE) up -d $(DB_SVC)

build:
	docker-compose -f $(DC_FILE) build

down:
	docker-compose -f $(DC_FILE) down

exec:
	docker-compose -f $(DC_FILE) exec $(SVC) sh

logs:
	docker logs -f --tail 100 $(PJT_NAME)_$(SVC)_1

dblogs:
	docker logs -f --tail 100 $(PJT_NAME)_$(DB_SVC)_1

rmvol:
	docker-compose -f $(DC_FILE) down -v