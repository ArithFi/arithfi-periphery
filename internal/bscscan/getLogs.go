package bscscan

import (
	"encoding/json"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs"
	"io"
	"net/http"
)

type Log struct {
	Address          string   `json:"address" bson:"address"`
	TimeStamp        string   `json:"timeStamp" bson:"timeStamp"`
	BlockNumber      string   `json:"blockNumber" bson:"blockNumber"`
	BlockHash        string   `json:"blockHash" bson:"blockHash"`
	GasPrice         string   `json:"gasPrice" bson:"gasPrice"`
	GasUsed          string   `json:"gasUsed" bson:"gasUsed"`
	LogIndex         string   `json:"logIndex" bson:"logIndex"`
	Data             string   `json:"data" bson:"data"`
	TransactionHash  string   `json:"transactionHash" bson:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex" bson:"transactionIndex"`
	Topics           []string `json:"topics" bson:"topics"`
}

type Result struct {
	Status  string `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
	Result  []Log  `json:"result" bson:"result"`
}

func GetLogs(fromBlock string, toBlock string) ([]Log, error) {
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
		return []Log{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []Log{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return []Log{}, err
	}

	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return []Log{}, err
	}

	return result.Result, nil
}
