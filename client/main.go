package client

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// CLIENT

// Client hols a connection to a specific server
type Client struct {
	connection *websocket.Conn
}

// NewClient returns a pointer to a new client node
func NewClient() *Client {
	currentSlave := Client{}
	return &currentSlave
}

// StartClientNode starts a client node in the specified host and port
func (c *Client) StartClientNode(host string, port int) {
	currentSlave := Client{}
	currentSlave.connectToServer(host, port)
	for {
		_, commandToRun, err := currentSlave.connection.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading from master: %+v\n", err)
		} else {
			fmt.Println(string(commandToRun))
			commandOutput := currentSlave.runShellCommand(commandToRun)
			err = currentSlave.connection.WriteMessage(2, commandOutput)
			if err != nil {
				fmt.Printf("Could not write message %s to ws %s\n", string(commandOutput), err)
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func (c *Client) connectToServer(host string, port int) {

	masterURL := fmt.Sprintf("ws://%s:%d/ws", host, port)
	connection, _, err := websocket.DefaultDialer.Dial(masterURL, nil)
	if err != nil {
		msg := fmt.Errorf("Error opening connection %+v", err)
		fmt.Println(msg)
	} else {
		c.connection = connection
	}

}

func formatCommand(command []byte) (string, []string) {
	stringCommand := string(command)
	splittedCommands := strings.Split(stringCommand, " ")
	var firstCommand string
	var arguments []string
	firstCommand = splittedCommands[0]

	if len(splittedCommands) > 1 {
		arguments = splittedCommands[1:]
	}

	return firstCommand, arguments
}

func (c *Client) runShellCommand(instructions []byte) []byte {
	command, arguments := formatCommand(instructions)

	log.Printf("Running %s %+v command\n", command, arguments)

	output, err := exec.Command(command, arguments...).Output()
	if err != nil {
		log.Printf("Error running command: %s\n", err)
	}
	fmt.Println(string(output))
	return output
}
