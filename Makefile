.PHONY: build
build:
	go build -v ./cmd/server

.PHONY: migrateup
migrateup:
	migrate -path migrations/ -database 'postgres://vbnm251:vbnm251@localhost:5432/sayrsa?sslmode=disable' up

.PHONY: migratedown
migratedown:
	migrate -path migrations/ -database 'postgres://c73442_sayrsa_na4u_ru:MiCfeDodqemos29@postgres.c73442.h2/c73442_sayrsa_na4u_ru' down

.DEFAULT_GOAL := build