include .env
export

up:
	docker-compose up -d

upv:
	docker-compose up

updb:
	docker-compose up -d db

down:
	docker-compose down

server:
	go run .

build:
	go build -o app

outenv:
	./app outenv

outenvfile:
	./app outenv > .env

new_migration:
	migrate create -ext sql -dir migration -seq $(name)

migrateup:
	migrate -path migration -database "$(DB_SOURCE)" -verbose up

migrateup1:
	migrate -path migration -database "$(DB_SOURCE)" -verbose up 1

migratedown:
	migrate -path migration -database "$(DB_SOURCE)" -verbose down

migratedown1:
	migrate -path migration -database "$(DB_SOURCE)" -verbose down 1

.PHONY: up upv updb down server build outenv outenvfile new_migration migrateup migratedown migrateup1 migratedown1