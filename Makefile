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

.PHONY: up upv updb down server build outenv outenvfile