package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	waSocket "github.com/amiruldev20/waSocket"
	waProto "github.com/amiruldev20/waSocket/binary/proto"
	"github.com/amiruldev20/waSocket/store"
	"github.com/amiruldev20/waSocket/store/sqlstore"
	"github.com/amiruldev20/waSocket/types/events"
	message "github.com/danzok/Alice-wabot/system"
	"github.com/joho/godotenv"
	"github.com/mdp/qrterminal"
	"github.com/probandula/figlet4go"
)

func main() {

	//load yours variable in .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//create you db sqlite3
	container,
		err := sqlstore.New("sqlite3", "file:yoursqlitefile.db?_foreign_keys=on", nil)
	if err != nil {
		log.Fatal(err)
	}

	// save the first device session
	deviceStore,
		err := container.GetFirstDevice()
	if err != nil {
		log.Fatal(err)
	}

	/* Setting env */
	typeLogin := os.Getenv("TYPE_LOGIN")
	numberBot := os.Getenv("NUMBER_BOT")

	//client//
	sock := waSocket.NewClient(deviceStore, nil)

	eventHandler := registerHandler(sock)
	sock.AddEventHandler(eventHandler)

	/* what type of login */

	if sock.Store.ID == nil {

		/*type code */
		if typeLogin == "code" {
			fmt.Println("You login with pairing code")
			fmt.Println("Bot Numer : " + numberBot)

			err = sock.Connect()
			if err != nil {
				log.Fatal(err)
			}

			/* don't edit */
			code,
				err := sock.PairPhone(numberBot, true, waSocket.PairClientChrome, "Chrome (Alice-bot)")

			if err != nil {
				log.Fatal(err)
			}
			log.Println("Your code : " + code)
		} else {
			/*type qr*/
			qrChan,
				_ := sock.GetQRChannel(context.Background())

			err = sock.Connect()
			if err != nil {
				panic(err)
			}
			for evt := range qrChan {
				if evt.Event == "code" {
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)

					log.Println("Please scan this QR...")
				} else {
					log.Println("Login successfully!! ✅")
				}
			}
		}

	} else {

		/* Already logged in, just connect */
		err = sock.Connect()
		log.Println("Login Sucessfully!!")
		if err != nil {
			panic(err)
		}
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	sock.Disconnect()
}

func init() {
	ascii := figlet4go.NewAsciiRender()

	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		// Colors can be given by default ansi color codes...
		figlet4go.ColorGreen,
		figlet4go.ColorYellow,
		figlet4go.ColorCyan,
		// ...or by an hex string...
		//figlet4go.NewTrueColorFromHexString("885DBA"),
		// ...or by an TrueColor object with rgb values
		//figlet4go.TrueColor{136, 93, 186},
	}

	renderStr, _ := ascii.RenderOpts("Alice beta", options)
	store.DeviceProps.PlatformType = waProto.DeviceProps_FIREFOX.Enum()
	//store.DeviceProps.Os = proto.String(string(dxz))
	fmt.Print(renderStr)
}

func registerHandler(sock *waSocket.Client) func(evt interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			if strings.HasPrefix(v.Info.ID, "BAE5") {
				return
			}

			//function to view contact status
			/*if v.Info.Chat.String() == "status@broadcast" {
				sock.MarkRead([]types.MessageID{v.Info.ID}, v.Info.Timestamp, v.Info.Chat, v.Info.Sender)
				//fmt.Println("")
			}*/
			go message.Msg(sock, v)
			break
		}
	}
}
