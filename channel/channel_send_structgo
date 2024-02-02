package main

import (
	"fmt"
	"time"
)

type Message struct {
	Text string
}

func main() {
	messageChan := make(chan Message, 1)

	go func() {
		messageChan <- Message{Text: "Hello, Channel!"}
	}()

	time.Sleep(1 * time.Second)
	receivedMessage := <-messageChan

	fmt.Println(receivedMessage.Text)
}
