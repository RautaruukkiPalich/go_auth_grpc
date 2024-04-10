run:
	go run cmd/auth/main.go --config=./config/local.yaml

dockerrun:
	docker-compose up -d --build 

migrate:
	go run ./cmd/migrator --config=./config/local.yaml

lint:
	golangci-lint run ./...

makemigrations:
	migrate create -ext sql -dir migrations $(name)