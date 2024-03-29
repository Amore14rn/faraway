<h1 align="center">Hi there, I'm <a href="https://github.com/Amore14rn" target="_blank">Roman</a> 

<h3 align="center">Faraway test task</h3>

- [Schema](#Schema)
- [AboutTask](#AboutTask)
- [Task](#Task)
- [Requirements](#Requirements)
- [Installation](#Installation)
- [Testing](#Testing)



## Schema

````
faraway - |
          |__cmd|
          |     |__client|
          |     |        |__main.go
          |     |__server|
          |              |__main.go
          |__internal|
          |          |__client|
          |          |        |__client.go
          |          |        |__client_test.go
          |          |__server|
          |                   |__server.go
          |                   |__server_test.go
          |__pkg|
          |     |__pow|
          |     |     |__pow.go
          |     |     |__pow_test.go
          |     |config|
          |     |      |__config.go  
          |     |clock|
          |     |     |__clock.go
          |     |
          |     |protocol|
          |     |        |__protocol.go
          |     |        |__protocol_test.go
          |     |
          |     |redis|
          |     |     |__redis.go
          |     |     |__memory.go
          |
          |_.gitignore
          |_go.mod
          |_client.Dockerfile
          |_server.Dockerfile    
          |_docker-compose.yml
          |_README.md
          |_Makefile
````

## AboutTask
Test task for Server Engineer

Design and implement “Word of Wisdom” tcp server.  
- TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.  
- The choice of the POW algorithm should be explained.  
- After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.  
- Docker file should be provided both for the server and for the client that solves the POW challenge

## Task

- SHA-256 (Bitcoin):
  :Преимущества: Широко известен, имеет большое сообщество и высокий уровень безопасности благодаря своей распространенности и долгому времени работы.
  :Недостатки: Майнинг на ASIC-устройствах может привести к централизации, так как оборудование ASIC дорого и не доступно каждому.
- Ethash (Ethereum):
  :Преимущества: Ориентирован на GPU, что способствует более равномерному распределению майнинга. Майнинг на стандартных видеокартах доступен широкому кругу людей.
  :Недостатки: Требует большого объема памяти, что может быть непрактичным для мобильных устройств или устройств с ограниченной памятью.
- Scrypt (Litecoin):
  :Преимущества: Также ориентирован на GPU и менее подвержен централизации через ASIC. Относительно более доступный майнинг для обычных пользователей.
  :Недостатки: Возможно, менее безопасен, чем более распространенные алгоритмы.
- RandomX (Monero):
  :Преимущества: Сложный для реализации на ASIC, что поддерживает более равномерное распределение майнинга. Майнинг может быть выполнен на обычных компьютерах.
  :Недостатки: Требует большого объема памяти и вычислительной мощности, что может быть недоступно для некоторых устройств.

## Requirements

### Файл server/main.go
- Выводит сообщение "start server" при старте сервера.
- Загружает конфигурацию из файла "config/config.json" с помощью функции config.Load(). Если загрузка конфигурации завершается с ошибкой, выводит сообщение об ошибке и завершает выполнение.
- Инициализирует контекст (ctx) и передает в него загруженную конфигурацию с помощью context.WithValue().
- Инициализирует и добавляет в контекст экземпляр clock.SystemClock{} для работы с временем.
- Инициализирует и добавляет в контекст экземпляр redis.RedisCache для работы с кэшем Redis.
- Инициализирует генератор случайных чисел с использованием текущего времени для случайного выбора цитат.
- Запускает серверное приложение с помощью функции server.Run(), передавая контекст и адрес. Если запуск сервера завершается с ошибкой, выводит сообщение об ошибке.

### Файл client/main.go
- Выводит сообщение "start client" при старте клиента.
- Загружает конфигурацию из файла "config/config.json" с помощью функции config.Load(). Если загрузка конфигурации завершается с ошибкой, выводит сообщение об ошибке и завершает выполнение.
- Инициализирует контекст (ctx) и передает в него загруженную конфигурацию с помощью context.WithValue().
- Создает строку адреса, объединяя значения хоста и порта из конфигурации.
- Запускает клиентское приложение с помощью функции client.Run(), передавая контекст и адрес. Если запуск клиента завершается с ошибкой, выводит сообщение об ошибке.

### Файл client.go:
- NewClient(conn net.Conn) *Client: Эта функция создает и возвращает новый экземпляр клиента Client, связанный с переданным соединением conn. Она инициализирует читатель (reader) и писатель (writer) для этого клиента.
- Close() error: Эта метод функции Client закрывает соединение клиента. Это важно для корректного завершения соединения после использования клиента.
- RequestChallenge() (*pow.HashcashData, error): Этот метод запроса вызывает сервер для получения вызова (challenge). Он отправляет соответствующее сообщение на сервер, читает и разбирает ответ, включая полученные данные вызова, и возвращает указатель на структуру HashcashData, представляющую полученный вызов.
- ComputeAndSendHashcash(ctx context.Context, hashcash *pow.HashcashData, conf *config.Config) error: Этот метод вычисляет хэш-наличие (hashcash) для проверки Proof of Work и отправляет его обратно на сервер. Он принимает контекст, данные хэш-наличия и конфигурацию, вычисляет хэш-наличие согласно максимальному количеству итераций в конфигурации, маршалирует его в JSON и отправляет обратно на сервер.
- GetResultQuote() (string, error): Этот метод получает и возвращает результат (цитату) от сервера после успешного завершения работы с хэш-наличием. Он читает сообщение от сервера и извлекает полезную нагрузку (payload) цитаты.
- Run(ctx context.Context, address string) error: Это главная функция, которая управляет работой клиента. Она устанавливает соединение с сервером, создает экземпляр клиента, выполняет цикл, отправляя запросы на сервер, вычисляя и отправляя хэш-наличие, и получая результаты. Она также обрабатывает контекст и таймер для управления временем между запросами.

### Файл server.go:
- Quotes: Слайс строк, содержащий несколько цитат для ответов на запросы клиента.
- ErrQuit: Ошибка, которая используется для завершения соединения с клиентом.
- Run(ctx context.Context, address string) error: Это главная функция, которая запускает сервер, прослушивая указанный адрес. Она создает слушатель (listener) для принятия входящих подключений. Затем она в бесконечном цикле принимает подключения и запускает для каждого нового подключения функцию handleConnection() в отдельной горутине.
- handleConnection(ctx context.Context, conn net.Conn): Эта функция обрабатывает взаимодействие с клиентом через установленное соединение conn. Она читает входящие сообщения от клиента с помощью bufio.NewReader, вызывает функцию ProcessRequest() для обработки запроса клиента и отправляет ответ обратно с помощью функции sendMsg().
- ProcessRequest(ctx context.Context, msgStr string, clientInfo string) (*protocol.Message, error): Эта функция обрабатывает запрос от клиента в зависимости от типа сообщения. В зависимости от заголовка сообщения (типа запроса), она вызывает соответствующую функцию для обработки.
- handleChallengeRequest(ctx context.Context, clientInfo string) (*protocol.Message, error): Эта функция обрабатывает запрос клиента на получение задания (challenge). Она создает новое задание (challenge) для клиента, добавляет случайное значение в кеш и отправляет сообщение с заданием обратно клиенту.
- handleResourceRequest(ctx context.Context, clientInfo string, payload string) (*protocol.Message, error): Эта функция обрабатывает запрос клиента на ресурс. Она проверяет корректность hashcash, включая проверку наличия корректного случайного значения в кеше, а также проверку на протухание hashcash. Если всё верно, она вычисляет hashcash и возвращает клиенту случайную цитату.
- sendMsg(msg protocol.Message, conn net.Conn) error: Эта функция отправляет сообщение обратно клиенту через соединение conn. Сообщение перед отправкой преобразуется в строку и завершается символом новой строки.

### Файл pow.go
- Определение структуры HashcashData, представляющей данные для создания хеш-наличия (hashcash).
- Метод Stringify(), который преобразует данные hashcash в строку.
- Функция sha256Hash(), вычисляющая хеш SHA-1 для заданной строки.
- Функция IsHashCorrect(), проверяющая, удовлетворяет ли хеш требованиям по количеству ведущих нулей.
- Метод ComputeHashcash(), который вычисляет hashcash методом перебора (брутфорса) до тех пор, пока не будет найдено подходящее значение хеша.


### Файл Clock.go
- Ссодержит определение простой структуры для системных часов (SystemClock), которая имеет метод Now(), возвращающий текущее время с использованием пакета time
- Этот код предоставляет абстракцию для получения текущего времени через системные часы, путем обращения к методу Now() структуры SystemClock

### Файл protocol.go:

#### Константы:
- Quit: Представляет тип сообщения, указывающий, что либо сервер, либо клиент хочет закрыть соединение.
- RequestChallenge: Указывает на сообщение от клиента к серверу, запрашивающее новый вызов от сервера.
- ResponseChallenge: Указывает на сообщение от сервера клиенту, содержащее вызов для клиента.
- RequestResource: Указывает на сообщение от клиента к серверу, содержащее решение вызова или запрос ресурса.
- ResponseResource: Указывает на сообщение от сервера клиенту, содержащее полезную информацию, если решение верное, или ошибку в противном случае.

#### Структура Message:
- Message: Представляет структуру сообщения, содержащую два поля:
- Header: Целое число, указывающее тип сообщения.
- Payload: Строка, содержащая полезную нагрузку сообщения, которая может быть JSON, цитатой или быть пустой.

#### Метод Stringify:
- Stringify(): Метод, ассоциированный с структурой Message, который преобразует сообщение в строковый формат для отправки через соединение TCP. Формат: "Header|Payload".

#### Функция ParseMessage:
- ParseMessage(): Функция, которая принимает строку в качестве входных данных и пытается разобрать её в структуру Message. Она разделяет входную строку на части с использованием символа "|", а затем проверяет количество частей. Если частей одна или две, она пытается разобрать заголовок. В случае успешного разбора она создает структуру Message и возвращает её, включая разобранный заголовок и необязательную полезную нагрузку.

### Файл protocol_test.go:
- B этом файле представлены тесты для различных сценариев разбора сообщений с использованием библиотеки testify/assert. Каждый тест проверяет разные аспекты функции ParseMessage и утверждает, что ожидаемые результаты соответствуют действительным. Эти тесты помогают убедиться, что функция ParseMessage работает правильно для различных случаев.

### Файл redis.go:
- InitRedisCache: Эта функция инициализирует экземпляр RedisCache и устанавливает соединение с Redis-сервером. Принимает контекст ctx, адрес host и порт port. Создает и конфигурирует клиент Redis. В конце функции, используя этот клиент, устанавливается тестовое значение в Redis для проверки соединения. Возвращает экземпляр RedisCache и ошибку (если есть).
- Add: Добавляет случайное значение с заданным временем жизни (в секундах) в кеш Redis. Принимает ключ key и время expiration. Использует метод Set клиента Redis для установки значения ключа.
- Get: Проверяет наличие ключа key в кеше Redis. Возвращает true, если ключ существует и не просрочен, и ошибку (если есть).
- Delete: Удаляет ключ key из кеша Redis, используя метод Del клиента Redis.

### Файл memory.go:
- MemoryCache: Создает экземпляр InMemoryCache, который представляет собой кеш в памяти. Используется для тестирования вместо Redis. Принимает в качестве параметра объект Clock, который используется для определения текущего времени (это позволяет легко заменить реальное время на тестовое).
- Add: Добавляет случайное значение с заданным временем жизни (в секундах) в кеш. Принимает ключ key и время expiration. Записывает внутренние данные с меткой времени и временем жизни.
- Get: Проверяет наличие ключа key в кеше. Возвращает true, если ключ существует и не просрочен, и ошибку (если есть).
- Delete: Удаляет ключ key из кеша.

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

```` bash
    .PHONY: .install-linter
.install-linter:
	### INSTALL GOLANGCI-LINT ###
	[ -f $(GOLANCI_LINT) ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) $(GOLANCI_LINT_VERSION)
````

```` bash
   .PHONY: lint
lint: .install-linter
	### RUN GOLANGCI-LINT ###
	$(GOLANGCI_LINT) run ./... --config=./.golangci.yml
````

```` bash
   .PHONY: lint-fast
lint-fast: .install-linter
	$(GOLANGCI_LINT) run ./... --fast --config=./.golangci.yml

````

## Testing

1.Run the tests:
```bash 
go test
```





