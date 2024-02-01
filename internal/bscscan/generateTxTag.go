package bscscan

import (
	"math/big"
)

type UserMap = map[string]string

func GenerateTxTag(from string, to string, amountETH *big.Float, userMap UserMap) string {
	fromNickname := "Unknown User" + from[:7]
	toNickname := "Unknown User" + to[:7]
	howMuch := amountETH.Text('f', 2)

	if userMap[from] != "" {
		fromNickname = userMap[from]
	}
	if userMap[to] != "" {
		toNickname = userMap[to]
	}

	return fromNickname + " 向 " + toNickname + " 转账 " + howMuch + " ATF"
}
