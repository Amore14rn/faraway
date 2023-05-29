install:
	go mod download

test:
	go clean --testcache
	go test ./...

start-server:
	go run app/cmd/server/main.go

start-client:
	go run app/cmd/client/main.go

start:
	docker-compose up --abort-on-container-exit --force-recreate --build server --build client

