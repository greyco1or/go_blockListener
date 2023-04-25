package main

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetContractAddress(client *ethclient.Client, txid common.Hash) string {
	receipt_tx, err := client.TransactionReceipt(context.Background(), txid)
	if err != nil {
		return ""
	}
	return receipt_tx.ContractAddress.Hex()
}

func GetRealGasUsed(client *ethclient.Client, txid common.Hash) uint64 {
	receipt_tx, err := client.TransactionReceipt(context.Background(), txid)
	if err != nil {
		return 0
	}
	return receipt_tx.GasUsed
}

func GetRealGasPrice(baseFee uint64, maxFeeCap uint64, maxTipCap uint64) *big.Int {
	if baseFee+maxTipCap > maxFeeCap {
		return new(big.Int).SetUint64(maxFeeCap)
	} else {
		return new(big.Int).SetUint64(baseFee + maxTipCap)
	}
}
