package main

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {

	// Connect to a server
	var connected = false
	var connection *nats.Conn
	for !connected {
		nc, err := nats.Connect("nats://localhost:4222")
		if err != nil {
			fmt.Println("Received a error")
			fmt.Println(err)
			time.Sleep(1 * time.Second)
		} else {
			connected = true
			connection = nc
		}
	}

	// Channel Subscriber
	ch := make(chan *nats.Msg, 64)
	_, err := connection.ChanSubscribe("*", ch)
	if err != nil {
		fmt.Println("Received a error")
		fmt.Println(err)
	}

	for {

		select {
		case msg := <-ch:
			//fmt.Printf("Subject =  %s\n", string(msg.Subject))
			//fmt.Printf("Message = %s\n", string(msg.Data))
			fmt.Println(string(msg.Data))
		}
	}
}
