package message

import (
	"fmt"

	"github.com/amiruldev20/waSocket"
	"github.com/amiruldev20/waSocket/types/events"
)

func Msg(sock *waSocket.Client, msg *events.Message) {
	fmt.Sprintln(sock, msg)

}
