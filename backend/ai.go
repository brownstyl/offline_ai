package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CallPythonAi(message string) (string, error) {
	//initialize a map that will be ecoded as json reply
	bodyMessage := map[string]string{
		"message": message,
	}
	//convert the big map into a json byte
	jsonBytes, err := json.Marshal(bodyMessage)
	if err != nil {
		return "", fmt.Errorf("failed to encode Json %v", err)
	}
	//wrap it with a byte new reader as in strings on form but this goes with json
	wrapBody := bytes.NewReader(jsonBytes)

	//create an outbon request to python
	req, err := http.NewRequest("POST", "http://localhost:5000/generate", wrapBody)
	if err != nil {
		return "", fmt.Errorf("Cannot create a New outbond Request this time. please try again %v", err)
	}

	//set header for python server
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	//create the sender button
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Couldn't send your Request this time. %v", err)
	}
	//close the client everytime to avoid memory leakage
	defer resp.Body.Close()

	//Read all incoming respose from python...
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("An error occured when reading and loading request %v", err)
	}

	//reading all the response gives you a json format so decode it
	//declare an empty map to hold the disposal of the json data
	aiResponse := map[string]string{}

	unfoldErr := json.Unmarshal(bodyByte, &aiResponse)
	if unfoldErr != nil {
		return "", fmt.Errorf("Oops! sorry an error occured during the unmarshal process %v", unfoldErr)
	}

	return aiResponse["reply"], nil
}
