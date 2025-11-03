package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	// Listen on localhost:3000. Change to ":3000" to accept external connections.
	ln, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	log.Println("listening on localhost:3000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			// Don't exit the whole program on a single accept error.
			log.Printf("accept error: %v", err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	remote := conn.RemoteAddr().String()
	log.Printf("connected: %s", remote)

	// Optional: set a read timeout for idle connections
	// conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

	reader := bufio.NewReader(conn)
	buf := make([]byte, 4096) // 4KB buffer

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			// Only print the bytes actually read
			msg := string(buf[:n])
			fmt.Printf("[%s] message: %q\n", remote, msg)
		}

		if err != nil {
			if err == io.EOF {
				log.Printf("client disconnected: %s", remote)
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("read timeout from %s", remote)
			} else {
				log.Printf("read error from %s: %v", remote, err)
			}
			return
		}
	}
}
