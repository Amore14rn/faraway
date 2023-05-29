package protocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	Exit              = iota // When exit, both sides (server or client) should terminate connection
	ChallengeRequest         // From client to server - asks for a new challenge from server
	ChallengeResponse        // From server to client - message contains a challenge for the client
	ResourceRequest          // From client to server - message contains the solved challenge
	ResourceResponse         // From server to client - message contains useful info if the solution is correct, or error if not
)

type Data struct {
	Header  int
	Payload string
}

func (d *Data) Encode() string {
	return fmt.Sprintf("%d|%s", d.Header, d.Payload)
}

func DecodeData(input string) (*Data, error) {
	input = strings.TrimSpace(input)
	var dataType int

	sections := strings.Split(input, "|")
	if len(sections) < 1 || len(sections) > 2 {
		return nil, errors.New("data doesn't match protocol")
	}

	dataType, err := strconv.Atoi(sections[0])
	if err != nil {
		return nil, errors.New("unable to parse code")
	}
	data := Data{
		Header: dataType,
	}

	if len(sections) == 2 {
		data.Payload = sections[1]
	}
	return &data, nil
}
