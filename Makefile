run-docker-compose:
	docker-compose up --build -d

remove-docker-compose:
	docker-compose down

run-populate:
	go run populate/main.go

run-all: run-docker-compose run-populate

install-cli:
	go build -o bin/cli.exe cli/main.go
