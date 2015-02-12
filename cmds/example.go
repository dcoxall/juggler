package main

import (
	"github.com/dcoxall/juggler/utils"
	"github.com/dcoxall/juggler"
	"fmt"
	"time"
	"os"
)

func main() {
	port := <-utils.FindAvailablePort()
	instance := juggler.NewInstance(port, "pong")
	fmt.Printf("Starting on port %d\n", port)
	_, err := instance.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "WTF!!! %s", err)
		os.Exit(1)
	}
	for {
		if utils.IsPortFree(port) {
			fmt.Printf("Port claims it is free\n")
		} else {
			fmt.Printf("Port claims it is taken... YAY!\n")
		}
		time.Sleep(1 * time.Second)
	}
}
