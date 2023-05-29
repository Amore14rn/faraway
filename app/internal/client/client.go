package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/Amore14rn/faraway/app/internal/pkg/config"
	"github.com/Amore14rn/faraway/app/internal/pkg/pow"
	"github.com/Amore14rn/faraway/app/internal/pkg/protocol"
)

func Run(ctx context.Context, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	fmt.Println("connected to", address)
	defer conn.Close()

	// client will send new request every 5 seconds endlessly
	for {
		message, err := HandleConnection(ctx, conn, conn)
		if err != nil {
			return err
		}
		fmt.Println("quote result:", message)
		time.Sleep(5 * time.Second)
	}
}

func readConnMsg(reader *bufio.Reader) (string, error) {
	return reader.ReadString('\n')
}

// sendMsg - send protocol message to connection
func sendMsg(msg protocol.Data, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.Encode())
	_, err := conn.Write([]byte(msgStr))
	return err
}

// HandleConnection - scenario for TCP-client
// 1. request challenge from server
// 2. compute hashcash to check Proof of Work
// 3. send hashcash solution back to server
// 4. get result quote from server
// readerConn and writerConn divided to more convenient mock on testing
func HandleConnection(ctx context.Context, readerConn io.Reader, writerConn io.Writer) (string, error) {
	reader := bufio.NewReader(readerConn)

	// 1. requesting challenge
	err := sendMsg(protocol.Data{
		Header: protocol.ChallengeRequest,
	}, writerConn)
	if err != nil {
		return "", fmt.Errorf("err send request: %w", err)
	}

	// reading and parsing response
	msgStr, err := readConnMsg(reader)
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}
	msg, err := protocol.DecodeData(msgStr)
	if err != nil {
		return "", fmt.Errorf("err parse msg: %w", err)
	}
	var hashcash pow.ProofOfWork
	err = json.Unmarshal([]byte(msg.Payload), &hashcash)
	if err != nil {
		return "", fmt.Errorf("err parse hashcash: %w", err)
	}
	fmt.Println("got hashcash:", hashcash)

	// 2. got challenge, compute hashcash
	cfg := config.GetConfig()
	hashcash, err = hashcash.ComputeProofOfWork(cfg.HashCash.MaxIterations)
	if err != nil {
		return "", fmt.Errorf("err compute hashcash: %w", err)
	}
	fmt.Println("hashcash computed:", hashcash)
	// marshal solution to json
	byteData, err := json.Marshal(hashcash)
	if err != nil {
		return "", fmt.Errorf("err marshal hashcash: %w", err)
	}

	// 3. send challenge solution back to server
	err = sendMsg(protocol.Data{
		Header:  protocol.ChallengeResponse,
		Payload: string(byteData),
	}, writerConn)
	if err != nil {
		return "", fmt.Errorf("err send request: %w", err)
	}
	fmt.Println("challenge sent to server")

	// 4. get result quote from server
	msgStr, err = readConnMsg(reader)
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}
	msg, err = protocol.DecodeData(msgStr)
	if err != nil {
		return "", fmt.Errorf("err parse msg: %w", err)
	}
	return msg.Payload, nil
}
