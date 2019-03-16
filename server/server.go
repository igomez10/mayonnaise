package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type masterNode struct {
	slaves    map[*websocket.Conn]bool // connected clients
	broadcast chan []byte              // broadcast channel
}

func main() {
	currentNode := masterNode{}
	currentNode.slaves = make(map[*websocket.Conn]bool)
	currentNode.broadcast = make(chan []byte)

	// Configure websocket route
	http.HandleFunc("/ws", currentNode.handleConnections)

	// Start listening for incoming chat messages
	go currentNode.readInputCommands()
	go currentNode.broadcastCommandsToRun()

	// Start the server on localhost port 8000 and log any errors
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Stopped Server: ", err)
	}
}

func (n *masterNode) readInputCommands() {

	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter commands to run:")
		command, err := inputReader.ReadBytes('\n')
		command = command[:len(command)-1]
		if err != nil {
			fmt.Println("Error reading command to send")
		}
		fmt.Printf("Read %s from input\n", command)
		time.Sleep(1 * time.Second)
		n.broadcast <- command
	}

}

func (n *masterNode) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	fmt.Println("Reveived Connection")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	(*n).slaves[ws] = true

	n.readInputCommands()

}

func (n *masterNode) broadcastCommandsToRun() {
	for {
		fmt.Println("Listening for commands to run:")
		// Grab the next message from the broadcast channel
		command := <-n.broadcast
		fmt.Printf("Received message: %+v\n", string(command))
		// Send it out to every client that is currently connected
		for slave := range n.slaves {
			err := slave.WriteMessage(1, command)
			if err != nil {
				log.Printf("error: %v", err)
				slave.Close()
				delete((*n).slaves, slave)
			}
			_, slaveOutput, err := slave.ReadMessage()
			if err != nil {
				log.Println("Error reading response from slave")
			}
			fmt.Printf("Running %s returned \n %s", string(command), slaveOutput)
		}
	}
}

func (n *masterNode) readFromConnection() {
	for {

	}
}
