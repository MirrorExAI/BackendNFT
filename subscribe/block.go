package subscribe

import (
	"context"
	"fmt"
	_ "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/contract"
	"log"
	"strings"
)

type SubscribeBlock struct {
	client *ethclient.Client
}

func NewSubscribeBlock(url string) (*SubscribeBlock, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &SubscribeBlock{
		client: client,
	}, nil

}
func (subBlock *SubscribeBlock) Subscribe() {
	headers := make(chan *types.Header)
	sub, err := subBlock.client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			block, err := subBlock.client.BlockByNumber(context.Background(), header.Number)
			if err != nil {
				log.Fatal(err)
			}
			for _, tx := range block.Transactions() {
				//log.Println("tx hash ", tx.Hash())
				ctx := context.Background()
				result, _, _ := subBlock.client.TransactionByHash(ctx, tx.Hash())
				if strings.EqualFold("0xdAC17F958D2ee523a2206206994597C13D831ec7", result.To().String()) {
					contractABI := GetContractABI()
					DecodeTransactionInputData(contractABI, tx.Data())
					//1. 判断是否存在这个地址
					//database.IsExit()
					//2.  判断调用的方法，如果是approval方法授权的地址是我们的，tg机器人通知，已授权平台，授权金额
					//3.  判断调用的方法，如果是approval方法授权的地址是我们的，tg机器人通知，已他人，授权金额
					//4， 判断调用的方法，如果是transfer方法，转移多少钱
				}
				if strings.EqualFold("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", result.To().String()) {
					log.Println("tx hash ", tx.Hash())
					log.Println("=============================================================")
					contractABI := GetContractABI()
					DecodeTransactionInputData(contractABI, tx.Data())
					log.Println("=============================================================")

				}

				//log.Println()
			}

			fmt.Println(block.Hash().Hex())      // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			fmt.Println(block.Number().Uint64()) // 3477413
			//fmt.Println(block.Time().Uint64())     // 1529525947
			fmt.Println(block.Nonce())             // 130524141876765836
			fmt.Println(len(block.Transactions())) // 7
		}
	}
}
func GetContractABI() *abi.ABI {
	contractABI, err := abi.JSON(strings.NewReader(string(contract.TetherTokenABI)))
	if err != nil {
		//log.Fatal(err)
	}
	return &contractABI
}

// refer
// https://github.com/ethereum/web3.py/blob/master/web3/contract.py#L435
func DecodeTransactionInputData(contractABI *abi.ABI, data []byte) {
	methodSigData := data[:4]
	inputsSigData := data[4:]
	method, err := contractABI.MethodById(methodSigData)
	if err != nil {
		log.Fatal(err)
	}
	inputsMap := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(inputsMap, inputsSigData); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(inputsMap)
	}

	log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Printf("Method Name: %s\n", method.Name)
	fmt.Printf("Method inputs: %v\n", inputsMap)

	log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>")

}
