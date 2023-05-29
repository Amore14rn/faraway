package server

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/Amore14rn/faraway/app/internal/pkg/config"
	"github.com/Amore14rn/faraway/app/internal/pkg/pow"
	"github.com/Amore14rn/faraway/app/internal/pkg/protocol"
)

// Quotes - const array of quotes to respond on client's request
var Quotes = []string{
	"The only limit to our realization of tomorrow will be our doubts of today.",

	"Your time is limited, don't waste it living someone else's life.",

	"You miss 100% of the shots you don't take.",

	"The future belongs to those who believe in the beauty of their dreams.",

	"The best preparation for tomorrow is doing your best today.",
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
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accept connection: %w", err)
		}
		// Handle connections in a new goroutine.
		go handleConnection(ctx, conn)
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
func ProcessRequest(ctx context.Context, msgStr string, clientInfo string) (*protocol.Data, error) {
	cfg := config.GetConfig()
	msg, err := protocol.DecodeData(msgStr)
	if err != nil {
		return nil, err
	}
	// switch by header of msg
	switch msg.Header {
	case protocol.Exit:
		return nil, ErrQuit
	case protocol.ChallengeResponse:
		fmt.Printf("client %s requests challenge\n", clientInfo)
		// create new challenge for client
		clock := ctx.Value("clock").(Clock)
		cache := ctx.Value("cache").(Cache)
		date := clock.Now()

		// add new created rand value to cache to check it later on RequestResource stage
		// with duration in seconds
		randValue := rand.Intn(100000)
		err := cache.Add(randValue, int64(cfg.HashCash.RequiredZerosCount))
		if err != nil {
			return nil, fmt.Errorf("err add rand to cache: %w", err)
		}

		hashcash := pow.ProofOfWork{
			Version:    1,
			ZerosCount: cfg.HashCash.RequiredZerosCount,
			Date:       date.Unix(),
			Resource:   clientInfo,
			Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", randValue))),
			Counter:    0,
		}
		hashcashMarshaled, err := json.Marshal(hashcash)
		if err != nil {
			return nil, fmt.Errorf("err marshal hashcash: %v", err)
		}
		msg := protocol.Data{
			Header:  protocol.ChallengeResponse,
			Payload: string(hashcashMarshaled),
		}
		return &msg, nil
	case protocol.ChallengeRequest:
		fmt.Printf("client %s requests resource with payload %s\n", clientInfo, msg.Payload)
		// parse client's solution
		var hashcash pow.ProofOfWork
		err := json.Unmarshal([]byte(msg.Payload), &hashcash)
		if err != nil {
			return nil, fmt.Errorf("err unmarshal hashcash: %w", err)
		}
		// validate hashcash params
		if hashcash.Resource != clientInfo {
			return nil, fmt.Errorf("invalid hashcash resource")
		}

		clock := ctx.Value("clock").(Clock)
		cache := ctx.Value("cache").(Cache)

		// decoding rand from base64 field in received client's hashcash
		randValueBytes, err := base64.StdEncoding.DecodeString(hashcash.Rand)
		if err != nil {
			return nil, fmt.Errorf("err decode rand: %w", err)
		}
		randValue, err := strconv.Atoi(string(randValueBytes))
		if err != nil {
			return nil, fmt.Errorf("err decode rand: %w", err)
		}

		// if rand exists in cache, it means, that hashcash is valid and really challenged by this server in past
		exists, err := cache.Get(randValue)
		if err != nil {
			return nil, fmt.Errorf("err get rand from cache: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("challenge expired or not sent")
		}

		// sent solution should not be outdated
		if clock.Now().Unix()-hashcash.Date > cfg.HashCash.ChallengeLifetime {
			return nil, fmt.Errorf("challenge expired")
		}
		//to prevent indefinite computing on server if client sent hashcash with 0 counter
		maxIter := hashcash.Counter
		if maxIter == 0 {
			maxIter = 1
		}
		_, err = hashcash.ComputeProofOfWork(maxIter)
		if err != nil {
			return nil, fmt.Errorf("invalid hashcash")
		}
		//get random quote
		fmt.Printf("client %s succesfully computed hashcash %s\n", clientInfo, msg.Payload)
		msg := protocol.Data{
			Header:  protocol.ResourceResponse,
			Payload: Quotes[rand.Intn(4)],
		}
		// delete rand from cache to prevent duplicated request with same hashcash value
		cache.Delete(randValue)
		return &msg, nil
	default:
		return nil, fmt.Errorf("unknown header")
	}
}

// sendMsg - send protocol message to connection
func sendMsg(msg protocol.Data, conn net.Conn) error {
	msgStr := fmt.Sprintf("%s\n", msg.Encode())
	_, err := conn.Write([]byte(msgStr))
	return err
}
