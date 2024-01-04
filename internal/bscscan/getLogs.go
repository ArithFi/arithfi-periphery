package main

import (
	"encoding/json"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs"
	"github.com/arithfi/arithfi-periphery/configs/mongo"
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

	//balances := make(map[string]*big.Float)

	//for _, v := range result.Result {
	//	timeStamp := new(big.Int)
	//	timeStamp.SetString(strings.TrimPrefix(v.TimeStamp, "0x"), 16)
	//	loc, _ := time.LoadLocation("Asia/Shanghai")
	//	date := time.Unix(timeStamp.Int64(), 0).In(loc).Format("2006-01-02")
	//	amountWei := new(big.Int)
	//	amountWei.SetString(strings.TrimPrefix(v.Data, "0x"), 16)
	//	amountEth := ConvertWeiToEth(amountWei)
	//	from := v.Topics[0]
	//	to := v.Topics[1]
	//
	//	if balances[from] == nil {
	//		balances[from] = new(big.Float)
	//	}
	//	beforeFromBalance := balances[from]
	//	balances[from].Sub(balances[from], amountEth)
	//	afterFromBalance := balances[from]
	//	if balances[to] == nil {
	//		balances[to] = new(big.Float)
	//	}
	//	beforeToBalance := balances[to]
	//	balances[to].Add(balances[to], amountEth)
	//	afterToBalance := balances[to]
	//
	//	// raw 为 原始结构的 v，用于备份，方便调试
	//	// 输出一个结构，存储到MongoDB
	//	// { raw: v, timeStamp: timeStamp.Int64(), abstract: { from: from, to: to, amount: amountEth }, before: { from: beforeFromBalance, to: beforeToBalance}, after: { from:  afterFromBalance, to: afterToBalance} } }
	//	fmt.Printf("Date: %s, From: %s, To: %s, Transaction Hash: %s, Amount: %s ETH\n", date, from, to, v.TransactionHash, amountEth.Text('f', 6))
	//}

	// 先插入到 MongoDB
	collection := mongo.MONGODB.Database("chain-bsc").Collection("transfer-logs")

}
