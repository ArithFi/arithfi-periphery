package bscscan

import (
	"math/big"
)

type UserMap = map[string]string

func GenerateTxTag(from string, to string, amountETH *big.Float, userMap UserMap) string {
	fromNickname := "用户" + from[:7]
	toNickname := "用户" + to[:7]
	howMuch := amountETH.Text('f', 2)

	if userMap[from] != "" {
		fromNickname = userMap[from]
	}
	if userMap[to] != "" {
		toNickname = userMap[to]
	}

	return fromNickname + " 向 " + toNickname + " 转账 " + howMuch + " ATF"
}
