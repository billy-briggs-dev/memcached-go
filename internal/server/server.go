package server

import (
	"bufio"
	"fmt"
	"memcached-go/internal/lexer"
	"net"
	"net/http"
	"time"

	"github.com/dgraph-io/ristretto"
)

var cache *ristretto.Cache

func Init() {
	var err error
	cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M)
		MaxCost:     1 << 30, // maximum cost of cache (1GB)
		BufferItems: 64,      // number of keys per Get buffer
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create cache: %v", err))
	}
}

func Start(port int) {
	go func() {
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
		http.ListenAndServe(":8080", nil) // Health check on port 8080
	}()

	address := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Printf("Listening on %s\n", address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

type cacheItem struct {
	Data  []byte
	Flags uint32
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for {
		cmd, err := lexer.ScanCommand(scanner, conn)
		if err != nil {
			fmt.Fprintf(conn, "ERROR\r\n")
			break
		}

		switch cmd.Name {
		case "set":
			expiry := time.Duration(cmd.Exptime) * time.Second
			cache.SetWithTTL(cmd.Key, cacheItem{Data: cmd.Data, Flags: cmd.Flags}, int64(len(cmd.Data)), expiry)
			cache.Wait()
			if !cmd.NoReply {
				fmt.Fprintf(conn, "STORED\r\n")
			}
		case "get":
			value, found := cache.Get(cmd.Key)
			if found {
				item, ok := value.(cacheItem)
				if !ok {
					fmt.Fprintf(conn, "SERVER_ERROR type assertion failed\r\n")
					continue
				}
				fmt.Fprintf(conn, "VALUE %s %d %d\r\n", cmd.Key, item.Flags, len(item.Data))
				conn.Write(item.Data)
				fmt.Fprintf(conn, "\r\nEND\r\n")
			} else {
				fmt.Fprintf(conn, "END\r\n")
			}
		default:
			fmt.Fprintf(conn, "ERROR\r\n")
		}
	}
}
