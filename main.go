package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"go-venice/apps"
	"go-venice/configs"
	"sync"
	"time"
)

func main() {

	config := configs.NewEnvConfig()
	baseRpc := apps.NewEvmRpc(config.RpcUrl, apps.BASE_SEPOLIA)

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

	unsigned, err := baseRpc.CreateRawTransaction(config.Address, config.Address, "0.00001", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(unsigned)

	signed, err := baseRpc.SingRawTransaction(unsigned, config.PrivateKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(signed)

	txHash, err := baseRpc.Broadcast(signed)

	if err != nil {
		panic(err)
	}
	fmt.Println(txHash)

	wg := sync.WaitGroup{}
	txChan := make(chan *types.Transaction)

	wg.Add(2)
	go func() {
		defer wg.Done()

		for range 10 {
			tx, err := baseRpc.GetTransaction(txHash)
			if tx != nil && err == nil {
				txChan <- tx
				break
			}
			time.Sleep(time.Second)
		}

		close(txChan)
	}()
	go func() {
		defer wg.Done()
		tx := <-txChan
		fmt.Println("txChan consume :", tx)
	}()

	wg.Wait()

}
