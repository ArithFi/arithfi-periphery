package main

import (
	"github.com/arithfi/arithfi-periphery/internal/scan"
	"time"
)

func main() {
	for {
		err := scan.FFutureTrading()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 10)
	}
}
