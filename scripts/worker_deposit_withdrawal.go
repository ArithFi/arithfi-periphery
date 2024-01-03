package main

import (
	"github.com/arithfi/arithfi-periphery/internal/scan"
	"time"
)

func main() {
	for {
		err := scan.DepositWithdrawal()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Minute)
	}
}
