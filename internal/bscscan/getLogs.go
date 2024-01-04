package bscscan

import (
	"encoding/json"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs"
	"io"
	"math/big"
	"net/http"
)

type Log struct {
	Address          string   `json:"address"`
	TimeStamp        string   `json:"timeStamp"`
	BlockNumber      string   `json:"blockNumber"`
	BlockHash        string   `json:"blockHash"`
	GasPrice         string   `json:"gasPrice"`
	GasUsed          string   `json:"gasUsed"`
	LogIndex         string   `json:"logIndex"`
	Data             string   `json:"data"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	Topics           []string `json:"topics"`
}

type Result struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []Log  `json:"result"`
}

// ConvertWeiToEth 将 wei 单位转换为 ETH 单位。
func ConvertWeiToEth(wei *big.Int) *big.Float {
	weiInEth := new(big.Float).SetInt(wei)
	ethValue := new(big.Float).Quo(weiInEth, big.NewFloat(1e18))
	return ethValue
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

	// 返回 result
	fmt.Println(result)

	return result.Result, nil
}
