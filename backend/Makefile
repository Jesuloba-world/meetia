DB_USER ?= postgres
DB_PASS ?= postgres
DB_NAME ?= meetia
DB_HOST ?= localhost
DB_PORT ?= 5432

DB_URL = postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: migrate-up migrate-down migrate-create migrate-status

.DEFAULT_GOAL := migrate-status

migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

migrate-create:
	goose -dir migrations create $(name) sql

migrate-status:
	goose -dir migrations postgres "$(DB_URL)" status

# make migrate-up DB_USER=[user] DB_PASS=[password] DB_HOST=[host] DB_PORT=[port] DB_NAME=[database]
