PG_PATH=infra/persistence/pg
PG_CONFIG_PATH=infra/persistence/sqlc.yml

include .env
export $(shell sed 's/=.*//' .env)

pg-code:
	sqlc generate -f ${PG_CONFIG_PATH}

pg-migrate:
	migrate create -ext sql -dir ${PG_PATH}/migration -seq $(name)

pg-make-migration:
	migrate -path ${PG_PATH}/migration -verbose -database "${PG_URL}" up $(n)

pg-rollback:
	migrate -path ${PG_PATH}/migration -verbose -database "${PG_URL}" down $(n)

pg-seed:
	echo 1

run:
	go run cmd/main.go
