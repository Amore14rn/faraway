<h1 align="center">Hi there, I'm <a href="https://github.com/Amore14rn" target="_blank">Roman</a> 

<h3 align="center">Faraway test task</h3>

- [AboutTask](#AboutTask)
- [Requirements](#Requirements)
- [Installation](#Installation)
- [Testing](#Testing)


faraway - |
          |__cmd|
                |__client|
                |        |__main.go
                |__server|
                         |__main.go
          |__internal|
                     |__client|
                     |         |__client.go
                     |         |__client_test.go
                     |__server|
                              |__server.go
                              |__server_test.go
          |__pkg|
                |__pow|
                |      |__pow.go
                |      |__pow_test.go
                |config|
                |       |__config.go  
                |clock|
                |     |__clock.go
                |
                |protocol|
                |        |__protocol.go
                |        |__protocol_test.go
                |
                |redis|
                |     |__redis.go
                |     |__memory.go
          |
          |_.gitignore
          |_go.mod
          |_client.Dockerfile
          |_server.Dockerfile    
          |_docker-compose.yml
          |_README.md
          |_Makefile


## AboutTask
Test task for Server Engineer

Design and implement “Word of Wisdom” tcp server.  
- TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.  
- The choice of the POW algorithm should be explained.  
- After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.  
- Docker file should be provided both for the server and for the client that solves the POW challenge

## Requirements
- 
-

## Installation
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


## Testing

1.Run the tests:
```bash 
go test
```


