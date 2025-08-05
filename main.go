package main

import (
	"fmt"
	"go-venice/apps"
	"go-venice/configs"
)

func main() {

	config := configs.NewEnvConfig()
	baseRpc := apps.NewEvmRpc(config.RpcUrl, apps.BASE_MAINNET)

	blockNumber, err := baseRpc.GetBlockNumber()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(blockNumber)

	res, err := baseRpc.GetBalance(config.Address)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("GetBalance: ", res)

	if err != nil {
		fmt.Println("Error converting hex to int:", err)
		return
	}

	delegated, err := baseRpc.GetDelegated(config.Address)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(delegated)

	reward, err := baseRpc.GetReward(config.Address)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reward)

	unsigned, err := baseRpc.CreateClaimTransaction(config.Address)
	fmt.Println(unsigned)
	if err != nil {
		panic(err)
	}

	signed, err := baseRpc.SingRawTransaction(unsigned, config.PrivateKey)
	fmt.Println(signed)
	if err != nil {
		panic(err)
	}
	/*
		tx, err := baseRpc.Broadcast(signed)
		fmt.Println(err)
		fmt.Println(tx)
	*/
}
