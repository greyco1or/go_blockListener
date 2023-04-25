package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// web3.js에서 WebsocketProvider가 ethclient.Dial 기반
func BlockListener() error {
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/246a412f52f74832b07b645d3c9b9fed")
	if err != nil {
		log.Fatal(err)
	}

	headers := make(chan *types.Header)
	//헤더정보가 새로 생겼을 때 받아온다.
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("####################Block Header#########################")
			fmt.Println("BLock Hash : ", block.Hash().Hex())
			fmt.Println("Block Number : ", block.Number().Uint64())
			fmt.Println("Block Timestamp : ", block.Time())
			fmt.Println("Block Nonce : ", block.Nonce())
			//블록에 포함된 트랜잭션 개수
			fmt.Println("Total Transaction : ", len(block.Transactions()))
			fmt.Println("####################Block Header#########################")
			baseFee := block.BaseFee()

			if len(block.Uncles()) > 0 {
				for _, uncle := range block.Uncles() {
					fmt.Println("-----This is Uncle Block Info Start----")
					//uncle블록의 수수로 계산법은 다르다.
					uncleFee := float64((uncle.Number.Uint64()+8-block.Number().Uint64())*2) / 8.0
					fmt.Println("Uncle Block Length : ", len(block.Uncles()))
					//블록 채굴자와 uncle 블록 채굴자는 다른 경우가 많다.
					fmt.Println("Uncle Miner Address : ", uncle.Coinbase.Hex())
					fmt.Println("Uncle Block Number : ", uncle.Number.Uint64())
					fmt.Println("Uncle Block Reward : ", uncleFee)
					fmt.Println("-----This is Uncle Block Info End----")
				}
			}

			for _, tx := range block.Transactions() {
				fmt.Println("*******************Transaction Info****************************")
				fmt.Println("Transaction Hash : ", tx.Hash().Hex())
				if tx.To() != nil {
					fmt.Println("To Address: ", tx.To())
				} else {
					fmt.Println(("To Address : Contract Creation"))
					contractAddress := GetContractAddress(client, tx.Hash())
					fmt.Println("Contract Address : ", contractAddress)
				}
				fmt.Println("Transfer Value(wei) : " + tx.Value().String())
				fmt.Println("Transaction nonce : ", tx.Nonce())
				fmt.Println("Transaction Gas Limit : ", tx.Gas()) //예상 Gas Limit

				realGasLimit := GetRealGasUsed(client, tx.Hash())
				fmt.Println("Transaction Real Gas Limit : ", realGasLimit)
				fmt.Println("Transaction GasFeeCap : ", tx.GasFeeCap().Uint64())
				fmt.Println("Transaction GasTipCap : ", tx.GasTipCap().Uint64())
				realGasPrice := GetRealGasPrice(baseFee.Uint64(), tx.GasFeeCap().Uint64(), tx.GasTipCap().Uint64())
				fmt.Println("Transaction RealGasPrice : ", realGasPrice)
				fmt.Println("Transaction Input Data : ", hex.EncodeToString(tx.Data()))
				fmt.Println("*******************Transaction Info****************************")

				//data가 있으면 컨트랙트의 함수를 호출하는 경우가 대부분
				if len(tx.Data()) != 0 {
					to, value := ERC20Transaction(hex.EncodeToString(tx.Data()))
					if to != "" {
						symbol, name, decimal := GetContractInfo(client, tx.To())
						fmt.Println("ERC20 Contract Address : ", tx.To().Hex())
						fmt.Println("ERC20 Contract Name : ", name)
						fmt.Println("ERC20 Contract Symbol : ", symbol)
						fmt.Println("ERC20 Contract Decimal : ", decimal)
						fmt.Println("ERC20 Transfer To Address : ", to)
						fmt.Println("ERC20 Transfer To Value : ", value) //value를 10**decimal로 나누어야한다.
					}
				}
			}

		}
	}
}
