package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/igomez10/mayonnaise/client"
	"github.com/igomez10/mayonnaise/server"
)

func TestServerRegistration(t *testing.T) {

	host := "localhost"
	port := 8000

	myserver := server.NewMasterNode()
	go myserver.StartMasterNode(host, port)

	myclient := client.NewClient()
	go myclient.StartClientNode(host, port)
	time.Sleep(100 * time.Millisecond)
	if len(myserver.Slaves) != 1 {
		t.Errorf("Number of connections is not matching, expected %d got %d", 1, len(myserver.Slaves))
	}
	fmt.Println()
}
