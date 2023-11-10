package lib

import (
	"context"
	"fmt"
	"strings"

	"github.com/amiruldev20/waSocket"
	waProto "github.com/amiruldev20/waSocket/binary/proto"
	"github.com/amiruldev20/waSocket/types"
	"github.com/amiruldev20/waSocket/types/events"
	"google.golang.org/protobuf/proto"
)

type alice struct {
	sock *waSocket.Client
	Msg  *events.Message
}

func NewSimp(Cli *waSocket.Client, m *events.Message) *alice {
	return &alice{
		sock: Cli,
		Msg:  m,
	}
}

/* parse jid */
func (m *alice) parseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient,
			err := types.ParseJID(arg)
		if err != nil {
			fmt.Printf("Invalid JID %s: %v", arg, err)
			return recipient, false
		} else if recipient.User == "" {
			fmt.Printf("Invalid JID %s: no server specified", arg)
			return recipient, false
		}
		return recipient,
			true
	}
}

/* send react */
func (m *alice) React(react string) {
	_,
		err := m.sock.SendMessage(context.Background(), m.Msg.Info.Chat, m.sock.BuildReaction(m.Msg.Info.Chat, m.Msg.Info.Sender, m.Msg.Info.ID, react))
	if err != nil {
		return
	}
}

/* send message */
func (m *alice) SendMsg(jid types.JID, teks string) {
	_,
		err := m.sock.SendMessage(context.Background(), jid, &waProto.Message{Conversation: proto.String(teks)})
	if err != nil {
		return
	}
}

/* send reply */
func (m *alice) Reply(teks string) {
	_,
		err := m.sock.SendMessage(context.Background(), m.Msg.Info.Chat, &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(teks),
			ContextInfo: &waProto.ContextInfo{
				Expiration:    proto.Uint32(86400),
				StanzaId:      &m.Msg.Info.ID,
				Participant:   proto.String(m.Msg.Info.Sender.String()),
				QuotedMessage: m.Msg.Message,
			},
		},
	})
	if err != nil {
		return
	}
}

/* send adReply */
func (m *alice) ReplyAd(teks string) {
	var isImage = waProto.ContextInfo_ExternalAdReplyInfo_IMAGE
	_, err := m.sock.SendMessage(context.Background(), m.Msg.Info.Chat, &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(teks),
			ContextInfo: &waProto.ContextInfo{
				ExternalAdReply: &waProto.ContextInfo_ExternalAdReplyInfo{
					Title:                 proto.String(""),
					Body:                  proto.String(""),
					MediaType:             &isImage,
					ThumbnailUrl:          proto.String(""),
					MediaUrl:              proto.String(""),
					SourceUrl:             proto.String("S"),
					ShowAdAttribution:     proto.Bool(true),
					RenderLargerThumbnail: proto.Bool(true),
				},
				Expiration:    proto.Uint32(86400),
				StanzaId:      &m.Msg.Info.ID,
				Participant:   proto.String(m.Msg.Info.Sender.String()),
				QuotedMessage: m.Msg.Message,
			},
		},
	})
	if err != nil {
		return
	}
}

/* send contact */
func (m *alice) SendContact(jid types.JID, number string, nama string) {
	_,
		err := m.sock.SendMessage(context.Background(), jid, &waProto.Message{
		ContactMessage: &waProto.ContactMessage{
			DisplayName: proto.String(nama),
			Vcard:       proto.String(fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nN:%s;;;\nFN:%s\nitem1.TEL;waid=%s:+%s\nitem1.X-ABLabel:Mobile\nEND:VCARD", nama, nama, number, number)),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &m.Msg.Info.ID,
				Participant:   proto.String(m.Msg.Info.Sender.String()),
				QuotedMessage: m.Msg.Message,
			},
		},
	})
	if err != nil {
		return
	}
}

/* create channel */
func (m *alice) createChannel(params []string) {
	_,
		err := m.sock.CreateNewsletter(waSocket.CreateNewsletterParams{
		Name: strings.Join(params, " "),
	})
	if err != nil {
		return
	}
}

/* fetch group admin */
func (m *alice) FetchGroupAdmin(Jid types.JID) ([]string, error) {
	var Admin []string
	resp, err := m.sock.GetGroupInfo(Jid)
	if err != nil {
		return Admin, err
	} else {
		for _, group := range resp.Participants {
			if group.IsAdmin || group.IsSuperAdmin {
				Admin = append(Admin, group.JID.String())
			}
		}
	}
	return Admin, nil
}

/* get group admin */
func (m *alice) GetGroupAdmin(jid types.JID, sender string) bool {
	if !m.Msg.Info.IsGroup {
		return false
	}
	admin, err := m.FetchGroupAdmin(jid)
	if err != nil {
		return false
	}
	for _, v := range admin {
		if v == sender {
			return true
		}
	}
	return false
}

/* get link group */
func (m *alice) LinkGc(Jid types.JID, reset bool) string {
	link,
		err := m.sock.GetGroupInviteLink(Jid, reset)

	if err != nil {
		panic(err)
	}
	return link
}

func (m *alice) GetCMD() string {
	extended := m.Msg.Message.GetExtendedTextMessage().GetText()
	text := m.Msg.Message.GetConversation()
	imageMatch := m.Msg.Message.GetImageMessage().GetCaption()
	videoMatch := m.Msg.Message.GetVideoMessage().GetCaption()
	//pollVote := m.Msg.Message.GetPollUpdateMessage().GetVote()
	tempBtnId := m.Msg.Message.GetTemplateButtonReplyMessage().GetSelectedId()
	btnId := m.Msg.Message.GetButtonsResponseMessage().GetSelectedButtonId()
	listId := m.Msg.Message.GetListResponseMessage().GetSingleSelectReply().GetSelectedRowId()
	var command string
	if text != "" {
		command = text
	} else if imageMatch != "" {
		command = imageMatch
	} else if videoMatch != "" {
		command = videoMatch
	} else if extended != "" {
		command = extended
		/*
		   } else if pollVote != "" {
		   command = pollVote
		*/
	} else if tempBtnId != "" {
		command = tempBtnId
	} else if btnId != "" {
		command = btnId
	} else if listId != "" {
		command = listId
	}
	return command
}
