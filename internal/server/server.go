package server

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Amore14rn/faraway/pkg/config"
	"github.com/Amore14rn/faraway/pkg/pow"
	"github.com/Amore14rn/faraway/pkg/protocol"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// Quotes - const array of quotes to respond on client's request
var Quotes = []string{
	"Wisdom is the daughter of experience.",

	"Knowledge speaks, but wisdom listens.",

	"When given the choice, be wise in actions over wise in words.",

	"The wise don't wait to enjoy the view from the top of the mountain, they recognize the beauty in the journey upwards",
}

var ErrQuit = errors.New("client requests to close connection")

// Clock  - interface for easier mock time.Now in tests
type Clock interface {
	Now() time.Time
}

// Cache - interface for add, delete and check existence of rand values for hashcash
type Cache interface {
	// Add - add rand value with expiration (in seconds) to cache
	Add(int, int64) error
	// Get - check existence of int key in cache
	Get(int) (bool, error)
	// Delete - delete key from cache
	Delete(int)
}

// Run - main function, launches server to listen on given address and handle new connections
func Run(ctx context.Context, address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	// Close the listener when the application closes.
	defer listener.Close()
	fmt.Println("listening", listener.Addr())

	// Создаем канал для таймера
	timerChan := time.NewTimer(1 * time.Second).C

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timerChan:
			// Listen for an incoming connection.
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("error accept connection:", err)
				continue
			}
			// Handle connections in a new goroutine.
			go handleConnection(ctx, conn)
			// Сброс таймера для следующей итерации
			timerChan = time.NewTimer(1 * time.Second).C
		}
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	fmt.Println("new client:", conn.RemoteAddr())
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("err read connection:", err)
			return
		}
		msg, err := ProcessRequest(ctx, req, conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("err process request:", err)
			return
		}
		if msg != nil {
			err := sendMsg(*msg, conn)
			if err != nil {
				fmt.Println("err send message:", err)
			}
		}
	}
}

// ProcessRequest - process request from client
// returns not-nil pointer to Message if needed to send it back to client
func ProcessRequest(ctx context.Context, msgStr string, clientInfo string) (*protocol.Message, error) {
	msg, err := protocol.ParseMessage(msgStr)
	if err != nil {
		return nil, err
	}

	switch msg.Header {
	case protocol.Quit:
		return nil, ErrQuit
	case protocol.RequestChallenge:
		return handleChallengeRequest(ctx, clientInfo)
	case protocol.RequestResource:
		return handleResourceRequest(ctx, clientInfo, msg.Payload)
	default:
		return nil, fmt.Errorf("unknown header")
	}
}

func handleChallengeRequest(ctx context.Context, clientInfo string) (*protocol.Message, error) {
	fmt.Printf("client %s requests challenge\n", clientInfo)

	conf := ctx.Value("config").(*config.Config)
	clock := ctx.Value("clock").(Clock)
	cache := ctx.Value("cache").(Cache)
	date := clock.Now()

	randValue := rand.Intn(100000)
	err := cache.Add(randValue, conf.HashcashDuration)
	if err != nil {
		return nil, fmt.Errorf("err add rand to cache: %w", err)
	}

	hashcash := pow.HashcashData{
		Version:    1,
		ZerosCount: conf.HashcashZerosCount,
		Date:       date.Unix(),
		Resource:   clientInfo,
		Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", randValue))),
		Counter:    0,
	}
	hashcashMarshaled, err := json.Marshal(hashcash)
	if err != nil {
		return nil, fmt.Errorf("err marshal hashcash: %v", err)
	}
	msg := protocol.Message{
		Header:  protocol.ResponseChallenge,
		Payload: string(hashcashMarshaled),
	}
	return &msg, nil
}

func handleResourceRequest(ctx context.Context, clientInfo string, payload string) (*protocol.Message, error) {
	fmt.Printf("client %s requests resource with payload %s\n", clientInfo, payload)

	var hashcash pow.HashcashData
	err := json.Unmarshal([]byte(payload), &hashcash)
	if err != nil {
		return nil, fmt.Errorf("err unmarshal hashcash: %w", err)
	}

	conf := ctx.Value("config").(*config.Config)
	clock := ctx.Value("clock").(Clock)
	cache := ctx.Value("cache").(Cache)

	randValueBytes, err := base64.StdEncoding.DecodeString(hashcash.Rand)
	if err != nil {
		return nil, fmt.Errorf("err decode rand: %w", err)
	}
	randValue, err := strconv.Atoi(string(randValueBytes))
	if err != nil {
		return nil, fmt.Errorf("err decode rand: %w", err)
	}

	exists, err := cache.Get(randValue)
	if err != nil {
		return nil, fmt.Errorf("err get rand from cache: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("challenge expired or not sent")
	}

	if clock.Now().Unix()-hashcash.Date > conf.HashcashDuration {
		return nil, fmt.Errorf("challenge expired")
	}
	maxIter := hashcash.Counter
	if maxIter == 0 {
		maxIter = 1
	}
	_, err = hashcash.ComputeHashcash(maxIter)
	if err != nil {
		return nil, fmt.Errorf("invalid hashcash")
	}

	fmt.Printf("client %s successfully computed hashcash %s\n", clientInfo, payload)
	msg := protocol.Message{
		Header:  protocol.ResponseResource,
		Payload: Quotes[rand.Intn(4)],
	}
	cache.Delete(randValue)
	return &msg, nil
}

// sendMsg - send protocol message to connection
func sendMsg(msg protocol.Message, conn net.Conn) error {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := conn.Write([]byte(msgStr))
	return err
}
