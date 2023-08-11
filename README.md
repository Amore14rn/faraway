##Faraway test

install:
```` bash
	go mod download
````
build:
```` bash
	go build -o bin/server app/cmd/server/main.go
	go build -o bin/client app/cmd/client/main.go
````
test:
```` bash
	go clean --testcache
	go test ./...
````

start-server:
```` bash
	go run app/cmd/server/main.go
````

start-client:
```` bash
	go run app/cmd/client/main.go
````
start:
```` bash
	docker-compose up --abort-on-container-exit --force-recreate --build server --build client
````



