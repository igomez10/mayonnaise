package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// CLIENT

func main() {
	connection, _, err := websocket.DefaultDialer.Dial("ws://localhost:8000/ws", nil)
	if err != nil {
		msg := fmt.Errorf("Error opening connection %+v", err)
		fmt.Println(msg)
	} else {
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Println("WS Connection Established with server")
		for {
			command, err := inputReader.ReadBytes('\n')
			command = command[:len(command)-1]
			if err != nil {
				fmt.Println("Error reading command to send")
			}
			err = connection.WriteMessage(2, command)
			if err != nil {
				fmt.Printf("Could not write message %s to ws %s\n", string(command), err)
			}

			_, commandOutput, err := connection.ReadMessage()
			if err != nil {
				fmt.Printf("Error from pong: %+v", err)
			} else {
				fmt.Println(string(commandOutput))
			}
			time.Sleep(time.Second * 1)
		}
	}
}
