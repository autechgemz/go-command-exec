package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
	Host     string   `yaml:"host"`
	Port     string   `yaml:"port"`
	Commands []string `yaml:"commands"`
}

type Config struct {
	Servers []ServerConfig `yaml:"servers"`
}

func main() {
	// load config file
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file: %s", err)
	}

	for _, server := range config.Servers {
		var authMethods []ssh.AuthMethod

		// connect to SSH agent
		agentSocket := os.Getenv("SSH_AUTH_SOCK")
		if agentSocket != "" {
			conn, err := net.Dial("unix", agentSocket)
			if err == nil {
				defer conn.Close()
				ag := agent.NewClient(conn)
				// Publickey method
				authMethods = append(authMethods, ssh.PublicKeysCallback(ag.Signers))
			} else {
				log.Printf("Failed to connect to SSH agent: %s", err)

				// Password method
				authMethods = append(authMethods, ssh.Password(server.Password))
			}
		} else {
			// Password method
			authMethods = append(authMethods, ssh.Password(server.Password))
		}

		sshConfig := &ssh.ClientConfig{
			User:            server.User,
			Auth:            authMethods,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		// SSH connection
		sshConn, err := ssh.Dial("tcp", server.Host+":"+server.Port, sshConfig)
		if err != nil {
			log.Printf("Failed to connect %s: %s", server.Host, err)
			continue
		}
		defer sshConn.Close()

		// execute commands
		for _, cmd := range server.Commands {
			// create session
			session, err := sshConn.NewSession()
			if err != nil {
				log.Printf("Failed to create session for %s: %s", server.Host, err)
				continue
			}

			// execute command
			start := time.Now()
			output, err := session.CombinedOutput(cmd)
			elapsed := time.Since(start)
			if err != nil {
				log.Printf("Failed to run command '%s' on %s: %s", cmd, server.Host, err)
				session.Close()
				continue
			}

			// print result
			fmt.Printf("Host: %s, ExecTime: %s, ExecCommand: %s\n%s\n", server.Host, elapsed, cmd, string(output))
			session.Close()
		}
	}
}
