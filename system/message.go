package message

import (
	"fmt"
	"log"
	"strings"

	"github.com/amiruldev20/waSocket"
	"github.com/amiruldev20/waSocket/types/events"
	"github.com/danzok/Alice-wabot/lib"
	"github.com/joho/godotenv"
)

var (
	prefix = "."
	self   = false
	owner  = "5591984155848"
)

func Msg(sock *waSocket.Client, msg *events.Message) {
	fmt.Sprintln(sock, msg)

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	m :=
		lib.NewSimp(sock, msg)

	//from := msg.Info.Chat
	sender := msg.Info.Sender.String()
	pushName := msg.Info.PushName
	isOwner := strings.Contains(sender, owner)
	//isAdmin := m.GetGroupAdmin(from, sender)
	//isBotAdm := m.GetGroupAdmin(from, botNumber + "@s.whatsapp.net")
	//isGroup := msg.Info.IsGroup
	args := strings.Split(m.GetCMD(), " ")
	command := strings.ToLower(args[0])
	//query := strings.Join(args[1: ], ` `)
	//extended := msg.Message.GetExtendedTextMessage()
	//quotedMsg := extended.GetContextInfo().GetQuotedMessage()
	//quotedImage := quotedMsg.GetImageMessage()
	//quotedVideo := quotedMsg.GetVideoMessage()
	//quotedSticker := quotedMsg.GetStickerMessage()

	// Self
	if self && !isOwner {
		return
	}

	//-- CONSOLE LOG
	fmt.Println("\n===============================\nNAME: " + pushName + "\nJID: " + sender + "\nTYPE: " + msg.Info.Type + "\nMessage: " + m.GetCMD() + "")
	//fmt.Println(m.Msg.Message.GetPollUpdateMessage().GetMetadata())

	if strings.EqualFold(m.GetCMD(), "bot") {
		m.Reply("Active")
	}

	switch command {
	case "oi":
		m.Reply("oi")
		m.React("😊")

	}
}
