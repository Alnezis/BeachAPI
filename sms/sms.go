package sms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Result struct {
	Status     string  `json:"status"`
	Code       int     `json:"code"`
	CallID     string  `json:"call_id"`
	Balance    float64 `json:"balance"`
	Cost       float64 `json:"cost"`
	StatusText string  `json:"status_text,omitempty"`
}

func SendSMS(phoneNumber string) *Result {

	client := &http.Client{}

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://sms.ru/code/call?phone=%s&api_id=27A9D540-7943-8F3B-6838-C5BDC477A61E", phoneNumber), nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return nil
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	var res Result
	err3 := json.Unmarshal(resp_body, &res)
	if err3 != nil {
		fmt.Println("whoops:", err3)
		//outputs: whoops: <nil>
	}

	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))

	return &res
}
