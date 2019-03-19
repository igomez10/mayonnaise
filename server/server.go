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
		fmt.Print(">>> ")
		command, err := inputReader.ReadBytes('\n')
		command = command[:len(command)-1]
		if err != nil {
			fmt.Println("Error reading command to send")
		}
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
		// Grab the next message from the broadcast channel
		command := <-n.broadcast
		// Send it out to every slave that is currently connected
		for slave := range n.slaves {
			n.writeToSocket(slave, command)
			fmt.Println(string(n.readFromSocket(slave)))
		}
	}
}

func (n *masterNode) readFromSocket(connection *websocket.Conn) []byte {
	_, slaveOutput, err := connection.ReadMessage()
	if err != nil {
		log.Println("Error reading response from slave")
	}
	return slaveOutput
}

func (n *masterNode) writeToSocket(connection *websocket.Conn, payload []byte) error {
	err := connection.WriteMessage(1, payload)
	if err != nil {
		log.Printf("error: %v", err)
		connection.Close()
		delete((*n).slaves, connection)
	}
	return err
}
