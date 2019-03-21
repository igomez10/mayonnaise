package main

import (
	"fmt"
	"time"

	"github.com/igomez10/mayonnaise/client"
	"github.com/igomez10/mayonnaise/server"
)

func main() {
	host := "localhost"
	port := 8000

	myserver := server.NewMasterNode()
	go myserver.StartMasterNode(host, port)

	myclient := client.NewClient()
	go myclient.StartClientNode(host, port)
	time.Sleep(100 * time.Millisecond)

	fmt.Println(len(myserver.Slaves) > 0)
}
