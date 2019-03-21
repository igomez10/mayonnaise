package server

import (
	"bufio"
	"bytes"
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

// MasterNode holds a map of current slaves and
// a broacast channel to send commands to
type MasterNode struct {
	Slaves    map[*websocket.Conn]bool // connected clients
	broadcast chan []byte              // broadcast channel
}

// NewMasterNode creates a new master node, initializes
// attributes and returns a pointer to the node
func NewMasterNode() *MasterNode {
	currentNode := MasterNode{}
	currentNode.Slaves = make(map[*websocket.Conn]bool)
	currentNode.broadcast = make(chan []byte)
	return &currentNode
}

// StartMasterNode starts a master node in the specified host and port
func (n *MasterNode) StartMasterNode(host string, port int) {
	// Configure websocket route
	http.HandleFunc("/ws", n.handleConnections)

	// Start listening for incoming chat messages
	go n.readInputCommands()
	go n.broadcastCommandsToRun()

	// Start the server on localhost port 8000 and log any errors
	address := fmt.Sprintf("%s:%d", host, port)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("Stopped Server: ", err)
	}
}

func (n *MasterNode) readInputCommands() {

	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter commands to run:")
		fmt.Print(">>> ")
		command, err := inputReader.ReadBytes('\n')
		command = bytes.Replace(command, []byte("\n"), nil, -1)
		if err != nil {
			fmt.Println("Error reading command to send")
		}
		time.Sleep(1 * time.Second)
		n.broadcast <- command
	}

}

func (n *MasterNode) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	(*n).Slaves[ws] = true
	fmt.Printf("connected client %d", &ws)
	n.readInputCommands()
}

func (n *MasterNode) broadcastCommandsToRun() {
	for {
		// Grab the next message from the broadcast channel
		command := <-n.broadcast
		// Send it out to every slave that is currently connected
		for slave := range n.Slaves {
			n.writeToSocket(slave, command)
			fmt.Println(string(n.readFromSocket(slave)))
		}
	}
}

func (n *MasterNode) readFromSocket(connection *websocket.Conn) []byte {
	_, slaveOutput, err := connection.ReadMessage()
	if err != nil {
		log.Println("Error reading response from slave")
	}
	return slaveOutput
}

func (n *MasterNode) writeToSocket(connection *websocket.Conn, payload []byte) error {
	err := connection.WriteMessage(1, payload)
	if err != nil {
		log.Printf("error: %v", err)
		connection.Close()
		delete((*n).Slaves, connection)
		fmt.Printf("disconnected client %+v", &connection)
	}
	return err
}
