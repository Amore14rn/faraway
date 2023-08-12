package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Amore14rn/faraway/pkg/config"
	"github.com/Amore14rn/faraway/pkg/pow"
	"github.com/Amore14rn/faraway/pkg/protocol"
	"io"
	"net"
	"time"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
	writer io.Writer
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: conn,
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) RequestChallenge() (*pow.HashcashData, error) {
	err := sendMsg(protocol.Message{Header: protocol.RequestChallenge}, c.writer)
	if err != nil {
		return nil, fmt.Errorf("err send request: %w", err)
	}

	msgStr, err := readConnMsg(c.reader)
	if err != nil {
		return nil, fmt.Errorf("err read msg: %w", err)
	}
	msg, err := protocol.ParseMessage(msgStr)
	if err != nil {
		return nil, fmt.Errorf("err parse msg: %w", err)
	}

	var hashcash pow.HashcashData
	err = json.Unmarshal([]byte(msg.Payload), &hashcash)
	if err != nil {
		return nil, fmt.Errorf("err parse hashcash: %w", err)
	}

	return &hashcash, nil
}

func (c *Client) ComputeAndSendHashcash(ctx context.Context, hashcash *pow.HashcashData, conf *config.Config) error {
	computedHashcash, err := hashcash.ComputeHashcash(conf.HashcashMaxIterations)
	if err != nil {
		return fmt.Errorf("err compute hashcash: %w", err)
	}

	byteData, err := json.Marshal(computedHashcash)
	if err != nil {
		return fmt.Errorf("err marshal hashcash: %w", err)
	}

	err = sendMsg(protocol.Message{Header: protocol.RequestResource, Payload: string(byteData)}, c.writer)
	if err != nil {
		return fmt.Errorf("err send request: %w", err)
	}

	return nil
}

func (c *Client) GetResultQuote() (string, error) {
	msgStr, err := readConnMsg(c.reader)
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}

	msg, err := protocol.ParseMessage(msgStr)
	if err != nil {
		return "", fmt.Errorf("err parse msg: %w", err)
	}

	return msg.Payload, nil
}

func Run(ctx context.Context, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	fmt.Println("connected to", address)
	defer conn.Close()

	client := NewClient(conn)

	// Создаем канал для таймера
	timerChan := time.NewTimer(0).C

	// client will send new request every 5 seconds endlessly
	for {
		select {
		case <-timerChan:
			hashcash, err := client.RequestChallenge()
			if err != nil {
				return err
			}
			fmt.Println("got hashcash:", hashcash)

			conf := ctx.Value("config").(*config.Config)
			err = client.ComputeAndSendHashcash(ctx, hashcash, conf)
			if err != nil {
				return err
			}
			fmt.Println("challenge sent to server")

			message, err := client.GetResultQuote()
			if err != nil {
				return err
			}
			fmt.Println("quote result:", message)

			// Сброс таймера для следующей итерации
			timerChan = time.After(5 * time.Second)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// readConnMsg - read string message from connection
func readConnMsg(reader *bufio.Reader) (string, error) {
	return reader.ReadString('\n')
}

// sendMsg - send protocol message to connection
func sendMsg(msg protocol.Message, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := conn.Write([]byte(msgStr))
	return err
}
