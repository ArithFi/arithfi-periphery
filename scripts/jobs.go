package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {

	url := "http://localhost:8080/fapi/update_daily_deposit_withdraw_metrics"
	method := "POST"

	payload := strings.NewReader(`{
    "date": "2023-12-01"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
