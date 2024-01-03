package main

import (
	"encoding/json"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs"
	"io"
	"math/big"
	"net/http"
	"strings"
)

type Result struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		Address         string `json:"address"`
		TimeStamp       string `json:"timeStamp"`
		BlockNumber     string `json:"blockNumber"`
		Data            string `json:"data"`
		TransactionHash string `json:"transactionHash"`
	}
}

// ConvertWeiToEth 将 wei 单位转换为 ETH 单位。
func ConvertWeiToEth(wei *big.Int) *big.Float {
	weiInEth := new(big.Float).SetInt(wei)
	ethValue := new(big.Float).Quo(weiInEth, big.NewFloat(1e18))
	return ethValue
}

func main() {
	fromBlock := "0"
	toBlock := "latest"
	tokenAddress := "0x00000000bA2ca30042001aBC545871380F570B1F"
	topic0 := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	url := "https://api.bscscan.com/api?module=logs&action=getLogs&fromBlock=" + fromBlock +
		"&toBlock=" + toBlock +
		"&address=" + tokenAddress +
		"&topic0=" + topic0 +
		"&apikey=" + configs.EnvBscScanAPI()
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
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

	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, v := range result.Result {
		amountWei := new(big.Int)
		amountWei.SetString(strings.TrimPrefix(v.Data, "0x"), 16)
		amountEth := ConvertWeiToEth(amountWei)
		fmt.Printf("Transaction Hash: %s, Amount: %s ETH\n", v.TransactionHash, amountEth.Text('f', 6))
	}
}
