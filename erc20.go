package main

import (
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/liyue201/erc20-go/erc20"
)

func ERC20Transaction(data string) (string, string) {
	//길이값이 일정한 것을 이용
	if len(data) != 136 {
		return "", "0"
	}
	methodId := data[:8]
	to := data[32:72]
	value := data[72:136]

	//erc20전송인지 아니지 체크
	if methodId != "a9059cbb" {
		return "", "0"
	}
	i := new(big.Int)

	valueStr := strings.TrimLeft(value, "0")
	i.SetString(valueStr, 16)
	return to, i.String()
}

func GetContractInfo(client *ethclient.Client, to *common.Address) (string, string, uint8) {
	//NewGGToken에 이미 erc20의 abi에 들어있기 때문에 to주소와 client정보만 넘기면 컨트랙트와 통신가능
	instance, err := erc20.NewGGToken(*to, client)
	if err != nil {
		log.Fatal(err)
	}
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	return name, symbol, decimals
}
