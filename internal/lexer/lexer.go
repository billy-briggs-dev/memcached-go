package lexer

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type Command struct {
	Name      string
	Key       string
	Flags     uint32
	Exptime   uint32
	ByteCount uint32
	NoReply   bool
	Data      []byte // Store data as []byte
}

func ScanCommand(r *bufio.Scanner, reader io.Reader) (*Command, error) {
	if !r.Scan() {
		return nil, errors.New("failed to read command line")
	}
	line := r.Text()
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil, errors.New("empty command")
	}

	cmdName := strings.ToLower(parts[0])
	switch cmdName {
	case "set", "cas", "add", "replace", "append", "prepend":
		// <cmd> <key> <flags> <exptime> <bytes> [noreply]
		// cas <key> <flags> <exptime> <bytes> <cas unique> [noreply]
		minArgs := 5
		if cmdName == "cas" {
			minArgs = 6
		}
		if len(parts) < minArgs {
			return nil, errors.New("invalid " + cmdName + " command format")
		}
		cmd := &Command{
			Name:    parts[0],
			Key:     parts[1],
			NoReply: false,
		}
		flags, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, errors.New("invalid flags")
		}
		cmd.Flags = uint32(flags)

		exptime, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, errors.New("invalid exptime")
		}
		cmd.Exptime = uint32(exptime)

		byteCount, err := strconv.Atoi(parts[4])
		if err != nil {
			return nil, errors.New("invalid byte count")
		}
		cmd.ByteCount = uint32(byteCount)

		argIdx := 5
		if cmdName == "cas" {
			if len(parts) < 6 {
				return nil, errors.New("missing cas unique value")
			}
			argIdx = 6
		}

		if len(parts) > argIdx && parts[argIdx] == "noreply" {
			cmd.NoReply = true
		}

		if cmd.ByteCount > 0 {
			data := make([]byte, cmd.ByteCount+2) // +2 for \r\n
			_, err := io.ReadFull(reader, data)
			if err != nil {
				return nil, errors.New("failed to read data block")
			}
			cmd.Data = data[:cmd.ByteCount]
		}

		return cmd, nil

	case "get":
		if len(parts) < 2 {
			return nil, errors.New("invalid get command format")
		}
		return &Command{
			Name: parts[0],
			Key:  parts[1],
		}, nil

	case "delete":
		// delete <key> [noreply]
		if len(parts) < 2 {
			return nil, errors.New("invalid delete command format")
		}
		cmd := &Command{
			Name:    parts[0],
			Key:     parts[1],
			NoReply: false,
		}
		if len(parts) > 2 && parts[2] == "noreply" {
			cmd.NoReply = true
		}
		return cmd, nil

	default:
		return nil, errors.New("unsupported command")
	}
}
