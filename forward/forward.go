package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"path/filepath"
	"github.com/pkg/sftp"
)

const (
	senderAddress   = "localhost:2221" // Sender server address
	receiverAddress = "localhost:2222" // Receiver server address
	sshUser         = "user"           // SFTP username
	privateKeyPath  = "../keys/id_rsa" // Path to private key
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run forward.go <file_path> <directory>")
		return
	}

	filePath := os.Args[1] 
	directory := os.Args[2] 

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

	// SSH client configuration
	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer), // Use the private key for authentication
		},
		// For testing only (accept any host key)
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), 
	}

	// Connect to the sender SFTP server
	senderConn, err := ssh.Dial("tcp", senderAddress, config)
	if err != nil {
		log.Fatalf("Failed to connect to sender: %v", err)
	}
	defer senderConn.Close()

	// Create an SFTP client for the sender
	senderClient, err := sftp.NewClient(senderConn)
	if err != nil {
		log.Fatalf("Failed to create SFTP client for sender: %v", err)
	}
	defer senderClient.Close()

	// Download the file from the sender
	remoteFile, err := senderClient.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open remote file: %v", err)
	}
	defer remoteFile.Close()

	// Connect to the receiver SFTP server
	receiverConn, err := ssh.Dial("tcp", receiverAddress, config)
	if err != nil {
		log.Fatalf("Failed to connect to receiver: %v", err)
	}
	defer receiverConn.Close()

	// Create an SFTP client for the receiver
	receiverClient, err := sftp.NewClient(receiverConn)
	if err != nil {
		log.Fatalf("Failed to create SFTP client for receiver: %v", err)
	}
	defer receiverClient.Close()

	// Upload the file to the receiver
	remoteFilePathReceiver := filepath.Join(directory, filepath.Base(filePath))
	receiverFile, err := receiverClient.Create(remoteFilePathReceiver)
	if err != nil {
		log.Fatalf("Failed to create remote file on receiver: %v", err)
	}
	defer receiverFile.Close()

	// Copy the file from the sender to the receiver
	_, err = io.Copy(receiverFile, remoteFile)
	if err != nil {
		log.Fatalf("Failed to upload file to receiver: %v", err)
	}

	fmt.Printf("File %s forwarded to %s\n", filePath, remoteFilePathReceiver)
}