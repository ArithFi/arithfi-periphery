package bscscan

import (
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

var UserTagMap = map[string]string{
	"0xdccbdbaee4d9d6639242f18f4eb08f4edad1a331": "ArithFi: System",
	"0x7c4fb3E5ba0a5D80658889715b307e66916f29b2": "ArithFi: Deployer",
	"0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38": "PancakeSwap",
	"0x0000000000000000000000000000000000000000": "NULL",
	"0xe26d976910D688083c8F9eCcB25e42345E5b95a0": "ArithFi: BSC-ETH-Bridge",
}

type Response struct {
	Tag string `json:"tag"`
}

func GetTag(c echo.Context) error {
	address := c.QueryParam("address")
	address = strings.ToLower(address)

	if UserTagMap[address] != "" {
		return c.JSON(http.StatusOK, Response{
			Tag: UserTagMap[address],
		})
	}

	db := mysql.ArithFiDB
	query, err := db.Query(`SELECT * FROM f_kol_info WHERE lower(walletAddress) = ?`, address)
	if err != nil {
		return err
	}
	if query.Next() {
		return c.JSON(http.StatusOK, Response{
			Tag: "KOL" + address[2:10],
		})
	}

	query2, err := db.Query(`SELECT * FROM f_user_assets WHERE lower(walletAddress) = ?`, address)
	if err != nil {
		return err
	}
	if query2.Next() {
		return c.JSON(http.StatusOK, Response{
			Tag: "User" + address[2:10],
		})
	}

	return c.JSON(http.StatusOK, Response{
		Tag: "Unknown User" + address[2:10],
	})
}
