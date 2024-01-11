package offchain

import (
	"github.com/arithfi/arithfi-periphery/configs/mysql"
)

type BlacklistUser struct {
	WalletAddress string `json:"walletAddress" bson:"walletAddress"`
	Type          int64  `json:"type" bson:"type"`
	Notes         string `json:"notes" bson:"notes"`
	TgName        string `json:"tgName" bson:"tgName"`
}

// GetUserBlacklist 获取用户黑名单
func GetUserBlacklist() ([]BlacklistUser, error) {
	query, err := mysql.MYSQL.Query(`SELECT walletAddress, type, notes, tgName
FROM f_user_blacklist
`)
	if err != nil {
		return []BlacklistUser{}, err
	}
	defer query.Close()
	var documents []BlacklistUser
	for query.Next() {
		var walletAddress string
		var typeInt int64
		var notes string
		var tgName string

		err := query.Scan(&walletAddress, &typeInt, &notes, &tgName)
		if err != nil {
			return []BlacklistUser{}, err
		}
		documents = append(documents, BlacklistUser{
			WalletAddress: walletAddress,
			Type:          typeInt,
			Notes:         notes,
			TgName:        tgName,
		})
	}
	if err != nil {
		return []BlacklistUser{}, err
	}
	return documents, nil
}
