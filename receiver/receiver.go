package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"

	"github.com/pkg/sftp"
)

const (
	sshPort        = ":2222"           // SFTP server port
	sshUser        = "user"            // SFTP username
	privateKeyPath = "../keys/id_rsa"  // Path to private key
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run receiver.go <upload_directory>")
		return
	}

	uploadDir := os.Args[1] 

	// Create the upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Read the private key
	key, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
	}

	// Parse the private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// SSH server configuration
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if conn.User() == sshUser {
				return nil, nil
			}
			return nil, fmt.Errorf("unknown user: %s", conn.User())
		},
	}
	config.AddHostKey(signer)

	// Start the SFTP server
	listener, err := net.Listen("tcp", sshPort)
	if err != nil {
		log.Fatalf("Failed to start SFTP server: %v", err)
	}
	defer listener.Close()
	fmt.Printf("SFTP server is listening on %s\n", sshPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn, config)
	}
}

func handleConnection(conn net.Conn, config *ssh.ServerConfig) {
	defer conn.Close()

	// Perform SSH handshake
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		log.Printf("Failed to handshake: %v", err)
		return
	}
	defer sshConn.Close()

	// Discard out-of-band requests
	go ssh.DiscardRequests(reqs)

	// Handle SFTP channels
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("Failed to accept channel: %v", err)
			continue
		}
        
		go func(in <-chan *ssh.Request) {
			for req := range in {
				req.Reply(req.Type == "subsystem" && string(req.Payload[4:]) == "sftp", nil)
			}
		}(requests)

		// Create SFTP server
		server, err := sftp.NewServer(channel)
		if err != nil {
			log.Printf("Failed to create SFTP server: %v", err)
			continue
		}
		defer server.Close()

		// Handle SFTP requests
		if err := server.Serve(); err != nil {
			log.Printf("SFTP server error: %v", err)
		}
	}
}