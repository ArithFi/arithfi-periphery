package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	// Start and end dates for the range
	startDate := time.Date(2023, 9, 25, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 9, 30, 0, 0, 0, 0, time.UTC)

	// Iterate over each day in the range
	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		url := "http://localhost:8080/fapi/update_daily_deposit_withdraw_metrics"
		//url := "http://localhost:8080/fapi/update_daily_burn_metrics"
		//url := "http://localhost:8080/fapi/update_daily_futures_metrics"
		//url := "http://localhost:8080/api/summary/update_daily_trade_metrics"
		method := "POST"

		// Create the payload with the current date
		payload := fmt.Sprintf(`{
		    "date": "%s"
		}`, date.Format("2006-01-02"))

		client := &http.Client{}
		req, err := http.NewRequest(method, url, bytes.NewBufferString(payload))

		if err != nil {
			fmt.Println(err)
			continue
		}
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Response for %s: %s\n", date.Format("2006-01-02"), string(body))
	}
}
