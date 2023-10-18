package tx

import (
	"context"
	"encoding/hex"
	"fmt"
	_ "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/contract"
	"log"
	"strings"
	"testing"
)

func Test_Subscribe(t *testing.T) {
	const rawUrl = "wss://mainnet.infura.io/ws/v3/dab126a4e1f444569c8f517a42cddda2"
	const txChCap = 30

	rpcClient, rpcErr := rpc.Dial(rawUrl)
	if rpcErr != nil {
		log.Fatal(rpcErr)
	}
	client := gethclient.New(rpcClient)

	txCh := make(chan common.Hash)
	sub, subErr := client.SubscribePendingTransactions(context.Background(), txCh)
	if subErr != nil {
		log.Fatal(subErr)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case tx := <-txCh:
			PrintTx(tx.String())
		case err := <-sub.Err():
			log.Fatal(err)
		}
	}
}

func ParseTransactionBaseInfo(tx *types.Transaction) {
	fmt.Printf("Hash: %s\n", tx.Hash().Hex())
	fmt.Printf("ChainId: %d\n", tx.ChainId())
	fmt.Printf("Value: %s\n", tx.Value().String())
	//fmt.Printf("From: %s\n", GetTransactionMessage(tx).From().Hex())
	fmt.Printf("To: %s\n", tx.To().Hex())
	fmt.Printf("Gas: %d\n", tx.Gas())
	fmt.Printf("Gas Price: %d\n", tx.GasPrice().Uint64())
	fmt.Printf("Nonce: %d\n", tx.Nonce())
	fmt.Printf("Transaction Data in hex: %s\n", hex.EncodeToString(tx.Data()))
}

func PrintTx(txhash string) {

	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/dab126a4e1f444569c8f517a42cddda2")
	if err != nil {
		//log.Fatal(err)
	}

	//txHash := "0xc3373524d0fb51a05241f2293e5a990e55c544b1b21743ee8bcc9c055e4afc3f"
	tx, isPending, err := client.TransactionByHash(context.Background(), common.HexToHash(txhash))
	if err != nil {
		//log.Fatal(err)
	}

	//fmt.Printf("tx isPending: %t\n", isPending)

	//

	if isPending {
		log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<NO BLOCK>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

		//ParseTransactionBaseInfo(tx)
		if strings.EqualFold(tx.To().Hex(), "0xdAC17F958D2ee523a2206206994597C13D831ec7") || strings.EqualFold(tx.To().Hex(), "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48") {
			log.Println("==============================================FROM CONTRACT============================================")
			log.Println("txhash", txhash)
			contractABI := GetContractABI()
			ParseTransactionBaseInfo(tx)
			DecodeTransactionInputData(contractABI, tx.Data())

			//transfer
			//approval

		}
	}
}
func GetContractABI() *abi.ABI {
	//rawABIResponse, err := GetContractRawABI(contractAddress, etherscanAPIKey)
	//if err != nil {
	//	log.Fatal(err)
	//}

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
