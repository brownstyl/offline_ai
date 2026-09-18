package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"offline_ai/backend/handlers" //import functions from handler (for handleFunc)
	"offline_ai/backend/services" //import functions from service (sendSms)

	"github.com/joho/godotenv" //.env
)

//server power house

func main() {
	//load from .env credential
	errow := godotenv.Load()
	if errow != nil {
		log.Fatalf("Couldn't read from .env %v", errow)
	}
	//retrieve all special credential to be used
	atUsername := os.Getenv("AT_USERNAME")
	atApiKey := os.Getenv("AT_APIKEY")
	atUrl := os.Getenv("AT_URL")

	//define the initial struct against what we have in our credential.
	send := &services.AfricasTalkingSender{
		Username: atUsername,
		Apikey:   atApiKey,
		Url:      atUrl,
	}

	//using servermux for routing all handlers
	mux := http.NewServeMux()

	//the handler control panel cordinator embaded with an annonymous function
	mux.HandleFunc("/incoming-sms", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleIncomingSms(w, r, send)
	})

	fmt.Println("server started Listening on Port http://localhost:8080")

	//server
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("An error occured and the server couldn't respond...")
		return
	}

}
