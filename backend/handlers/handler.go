package handlers

import (
	"fmt"
	"net/http"

	"offline_ai/backend/ai_handler"
	"offline_ai/backend/services"
)

func HandleIncomingSms(w http.ResponseWriter, r *http.Request, sender services.SmsSender) {
	//set constraints for routing
	if r.Method != http.MethodPost {
		//log errors if constraints are violated or not met
		http.Error(w, "Oops nothing here to see!!!", http.StatusNotFound)
		return
	}
	//collect incoming users request
	phone_number := r.FormValue("from")
	message_content := r.FormValue("text")

	//calling generate function to handle us our generated text.
	aiReply, aiErr := ai_handler.CallPythonAi(message_content)
	if aiErr != nil {
		fmt.Printf("An error occured while process your request to  the ai, please retry %v", aiErr)
		return
	}

	//send outgoing message/response back to the users.
	err := sender.Send(phone_number, aiReply)
	if err != nil {
		fmt.Println("Oops! Failed to send Message!!", err)
	}

	fmt.Println("=======Recieved Data========")
	fmt.Println(phone_number)
	fmt.Println(message_content)
}
