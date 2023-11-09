package main

import (
	"log"

	waSocket "github.com/amiruldev20/waSocket"
	"github.com/amiruldev20/waSocket/store/sqlstore"
	"github.com/joho/godotenv"
)

func main() {

	//load yours variable in .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//create you db sqlite3
	container, err := sqlstore.New("sqlite3", "file:yoursqlitefile.db?_foreign_keys=on", nil)
	if err != nil {
		log.Fatal(err)
	}

	// save the first device session
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Fatal(err)
	}

	sock := waSocket.NewClient(deviceStore, nil)

	//	eventHandler := registerHandler(sock)
}

/*
func registerHandler(sock *waSocket.Client) {

}*/
