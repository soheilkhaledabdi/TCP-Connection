package main

import (
	"log"
	"net"
)

// Message represents a message received from a client.
type Message struct {
	from    string
	payload []byte
}

// Server represents the TCP server.
type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
}

// NewServer creates a new Server instance.
func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}

// Start starts the TCP server and listens for connections.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
		return err
	}
	defer ln.Close()

	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)
	return nil
}

// acceptLoop accepts incoming connections and handles them.
func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		log.Printf("New connection from: %s", conn.RemoteAddr())

		go s.readLoop(conn)
	}
}

// readLoop reads data from the connection and processes it.
func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Read error from %s: %v", conn.RemoteAddr(), err)
			return
		}

		s.msgch <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}

		_, err = conn.Write([]byte("Thanks for the message"))
		if err != nil {
			log.Printf("Write error to %s: %v", conn.RemoteAddr(), err)
			return
		}
	}
}

func main() {
	server := NewServer(":3000")

	go func() {
		for msg := range server.msgch {
			log.Printf("Received message from %s: %s", msg.from, string(msg.payload))
		}
	}()

	log.Fatal(server.Start())
}
