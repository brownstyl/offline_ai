package main

import (
	"fmt"
	"net/http"
)

func HandleIncomingSms(w http.ResponseWriter, r *http.Request, sender SmsSender) {
	//set constraints for routing
	if r.Method != http.MethodPost {
		//log errors if constraints are violated or not met
		http.Error(w, "Oops nothing here to see!!!", http.StatusNotFound)
		return
	}
	phone_number := r.FormValue("from")
	message_content := r.FormValue("text")

	err := sender.Send(phone_number, "Message Recieved! Ai reply is coming soon...see you there")
	if err != nil {
		fmt.Println("Oops! Failed to send Message!!", err)
	}

	fmt.Println("=======Recieved Data========")
	fmt.Println(phone_number)
	fmt.Println(message_content)
}
