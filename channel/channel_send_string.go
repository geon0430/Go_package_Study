package main

import (
	"fmt"
)

func main() {
	messageChannel := make(chan string)

	go func() {
		messageChannel <- "Hello!"
	}()

	message := <-messageChannel

	fmt.Println(message)
}
