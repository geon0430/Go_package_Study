package main

import (
	"fmt"
	"time"
)

// Message 구조체는 데이터를 담고 있으며 채널을 통해 전달될 수 있습니다.
type Message struct {
	Text string
}

func main() {
	// Message 타입의 채널을 생성합니다.
	messageChan := make(chan Message, 1)

	// go 루틴을 사용하여 메시지를 채널에 보냅니다.
	go func() {
		messageChan <- Message{Text: "Hello, Channel!"}
	}()

	// 메인 루틴에서는 채널로부터 메시지를 받습니다.
	// 채널에서 메시지를 기다리는 동안 다른 작업을 수행할 수 있습니다. 여기에서는 시간 지연을 추가했습니다.
	time.Sleep(1 * time.Second)
	receivedMessage := <-messageChan

	// 받은 메시지를 출력합니다.
	fmt.Println(receivedMessage.Text)
}
