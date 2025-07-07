package config

import (
	"log"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func OpenSFTPConnection() (*sftp.Client, error) {
	// SFTP credentials from .env
	sftpHost := SFTP_HOST
	sftpPort := SFTP_PORT
	sftpUser := SFTP_USERNAME
	sftpPassword := SFTP_PASSWORD

	// SSH client configuration
	config := &ssh.ClientConfig{
		User: sftpUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sftpPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	// Connect to the SFTP server
	conn, err := ssh.Dial("tcp", sftpHost+":"+sftpPort, config)
	if err != nil {
		log.Fatalf("Failed to connect to SFTP server: %v", err)
	}

	// Create an SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatalf("Failed to create SFTP client: %v", err)
	}

	return client, nil
}
