package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type SmsSender interface {
	Send(to string, message string) error
}

type AfricasTalkingSender struct {
	username string
	apikey   string
	url      string
}

func (a *AfricasTalkingSender) Send(to string, message string) error {
	//set values applicable as a post for request
	formData := url.Values{}
	formData.Set("username", a.username)
	formData.Set("to", to)
	formData.Set("message", message)

	//replace all possible space with a (+)
	encodeForm := formData.Encode()

	//translating the encoded form into a string.
	bodyreader := strings.NewReader(encodeForm)

	//create the outbond request to their server...

	req, err := http.NewRequest("POST", a.url, bodyreader)
	if err != nil {
		return fmt.Errorf("Cannot establish a new request this time, try again. %v", err)
	}

	//at this junction we've gotta set the header to the server we're sending our request to.
	req.Header.Set("apikey", a.apikey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//send the request that you've to the destination you also want.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Could Send Message. %v", err)
	}

	//close the stream to avoid memory leakage
	defer resp.Body.Close()

	// log errors just if there's any failure
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("failed to send: invalid API key (Status %d)", resp.StatusCode)
		}
		if resp.StatusCode == http.StatusBadRequest {
			return fmt.Errorf("failed to send: bad request or invalid fields (Status %d)", resp.StatusCode)
		}
		return fmt.Errorf("failed to send: server returned status %d", resp.StatusCode)
	}

	fmt.Printf("Message Sent Successfully!!...Response Status %s\n", resp.Status)
	return nil

}
