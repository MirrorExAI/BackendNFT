package contract

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"

	_ "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"math/big"
	"strings"
	"testing"
)

func Test_SubscribeBlock(t *testing.T) {
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/dab126a4e1f444569c8f517a42cddda2")
	if err != nil {
		log.Fatal(err)
	}

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f

			block, err := client.BlockByNumber(context.Background(), header.Number)
			if err != nil {
				log.Fatal(err)
			}

			for _, tx := range block.Transactions() {

				//log.Println("tx hash ", tx.Hash())
				ctx := context.Background()
				result, _, _ := client.TransactionByHash(ctx, tx.Hash())
				if strings.EqualFold("0xdAC17F958D2ee523a2206206994597C13D831ec7", result.To().String()) {
					ParseTransactionBaseInfo(result)
					contractABI := GetContractABI()
					DecodeTransactionInputData(contractABI, tx.Data())
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
		fmt.Println("tx hash", txhash)
		//ParseTransactionBaseInfo(tx)
		if tx.To().String() == "0xdAC17F958D2ee523a2206206994597C13D831ec7" {
			log.Println("==============================================FROM CONTRACT============================================")
			//log.Println("txhash", txhash)
			contractABI := GetContractABI()
			DecodeTransactionInputData(contractABI, tx.Data())
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

func DecodeTransactionLogs(receipt *types.Receipt, contractABI *abi.ABI) {
	for _, vLog := range receipt.Logs {
		// topic[0] is the event name
		event, err := contractABI.EventByID(vLog.Topics[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Event Name: %s\n", event.Name)
		// topic[1:] is other indexed params in event
		if len(vLog.Topics) > 1 {
			for i, param := range vLog.Topics[1:] {
				fmt.Printf("Indexed params %d in hex: %s\n", i, param)
				fmt.Printf("Indexed params %d decoded %s\n", i, common.HexToAddress(param.Hex()))
			}
		}
		if len(vLog.Data) > 0 {
			fmt.Printf("Log Data in Hex: %s\n", hex.EncodeToString(vLog.Data))
			outputDataMap := make(map[string]interface{})
			err = contractABI.UnpackIntoMap(outputDataMap, event.Name, vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Event outputs: %v\n", outputDataMap)
		}
	}
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

func TestTetherTokenCaller_Allowance(t *testing.T) {

	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/dab126a4e1f444569c8f517a42cddda2")
	if err != nil {
		log.Fatal(err)
	}

	//txHash := "0xc3373524d0fb51a05241f2293e5a990e55c544b1b21743ee8bcc9c055e4afc3f"
	tx, isPending, err := client.TransactionByHash(context.Background(), common.HexToHash("0xc3373524d0fb51a05241f2293e5a990e55c544b1b21743ee8bcc9c055e4afc3f"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx isPending: %t\n", isPending)

	//log.Println(tx)
	if isPending {
		ParseTransactionBaseInfo(tx)
		contractABI := GetContractABI()
		//
		DecodeTransactionInputData(contractABI, tx.Data())
		//receipt := GetTransactionReceipt(client, txHash)
		//DecodeTransactionLogs(receipt, contractABI)
	}
}
func GetContractABI() *abi.ABI {
	//rawABIResponse, err := GetContractRawABI(contractAddress, etherscanAPIKey)
	//if err != nil {
	//	log.Fatal(err)
	//}

	contractABI, err := abi.JSON(strings.NewReader(string(TetherTokenABI)))
	if err != nil {
		//log.Fatal(err)
	}
	return &contractABI
}

func Test_SubscribeEvent(t *testing.T) {

	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/dab126a4e1f444569c8f517a42cddda2")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	//ABI := "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_upgradedAddress\",\"type\":\"address\"}],\"name\":\"deprecate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"deprecated\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_evilUser\",\"type\":\"address\"}],\"name\":\"addBlackList\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"upgradedAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maximumFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_maker\",\"type\":\"address\"}],\"name\":\"getBlackListStatus\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBasisPoints\",\"type\":\"uint256\"},{\"name\":\"newMaxFee\",\"type\":\"uint256\"}],\"name\":\"setParams\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"issue\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"basisPointsRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"isBlackListed\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_clearedUser\",\"type\":\"address\"}],\"name\":\"removeBlackList\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_UINT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_blackListedUser\",\"type\":\"address\"}],\"name\":\"destroyBlackFunds\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_initialSupply\",\"type\":\"uint256\"},{\"name\":\"_name\",\"type\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\"},{\"name\":\"_decimals\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Issue\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Redeem\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"Deprecate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"feeBasisPoints\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"maxFee\",\"type\":\"uint256\"}],\"name\":\"Params\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_blackListedUser\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_balance\",\"type\":\"uint256\"}],\"name\":\"DestroyedBlackFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"AddedBlackList\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"RemovedBlackList\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"}]"

	contractAbi, err := abi.JSON(strings.NewReader(string(TetherTokenABI)))

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	logTransferSig := []byte("Transfer(address,address,uint256)")
	LogApprovalSig := []byte("Approval(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash := crypto.Keccak256Hash(LogApprovalSig)
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:

			switch vLog.Topics[0].Hex() {
			case logTransferSigHash.Hex():
				fmt.Printf("Transfer Event")
				//fmt.Println(vLog) // pointer to event log
				from := common.HexToAddress(vLog.Topics[1].Hex())
				to := common.HexToAddress(vLog.Topics[2].Hex())
				log.Println(from)
				log.Println(to)

				event := make(map[string]interface{})
				// parse data
				if err := contractAbi.UnpackIntoMap(event, "Transfer", vLog.Data); err != nil {
					panic(err)
				}

				log.Println(event)

			case logApprovalSigHash.Hex():
				//fmt.Printf("Log Name: Approval\n")
				//from := common.HexToAddress(vLog.Topics[1].Hex())
				//to := common.HexToAddress(vLog.Topics[2].Hex())
				////value := vLog.Topics[0].Big().Int64()
				//log.Println(from)
				//log.Println(to)
				//
				//log.Println(vLog)
				//
				////var approvalEvent LogApproval
				////
				//
				//event := make(map[string]interface{})
				//// parse data
				//if err := contractAbi.UnpackIntoMap(event, "Approval", vLog.Data); err != nil {
				//	panic(err)
				//}
				//
				//log.Println(event)
				////err := contractAbi.Unpack(&approvalEvent, "Approval", vLog.Data)
				//err := contractAbi.UnpackIntoInterface(&approvalEvent, "Approval", vLog.Data)
				//if err != nil {
				//	log.Fatal(err)
				//}

				//log.Println("amount: ", approvalEvent.Tokens)

			}

		}
	}
}
func Test_Print(t *testing.T) {

	//var y *big.Int
	//
	//y.SetInt64(1000000)
	//
	//var x *big.Int
	//
	//result := x.Div(x.SetInt64(647742484), y.SetInt64(1000000))
	//
	//log.Println(result.Int64())

	var balance int64 = 2110000

	log.Println(float64(balance) / 1000000)

}

func Test_TransferAToB(t *testing.T) {
	conn, err := ethclient.Dial("https://mainnet.infura.io/v3/dab126a4e1f444569c8f517a42cddda2")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum web3client: %v", err)
	}
	store, err := NewTetherTokenTransactor(common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Storage contract: %v", err)
	}

	privateKey, err := crypto.HexToECDSA("56eaadc26990a5836ea890ff7a262703a4793fb60cfdb51af66c7ea17443b25b")
	if err != nil {
		log.Fatal(err)
	}

	nonce, _ := conn.NonceAt(context.Background(), common.HexToAddress("0xFF8C78F54235D04f88378f98f5ec6Fa68802f1b4"), nil)
	gasPrice, _ := conn.SuggestGasPrice(context.Background())
	//用哪条链，就用那个id
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))
	auth.GasLimit = uint64(300000)
	auth.Nonce = new(big.Int).SetUint64(nonce)
	auth.GasPrice = gasPrice
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	log.Println("Nonce: ", auth.Nonce)
	store2, err := NewTetherTokenCaller(common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Storage contract: %v", err)
	}
	amount, err := store2.BalanceOf(nil, common.HexToAddress("0xE3DB91F3B07F282f9890aE52e7fC0Bd2Adee2C74"))
	log.Println(amount)
	//// Call the store() function

	tx, err := store.TransferFrom(auth, common.HexToAddress("0xE3DB91F3B07F282f9890aE52e7fC0Bd2Adee2C74"), common.HexToAddress("0xFF8C78F54235D04f88378f98f5ec6Fa68802f1b4"), big.NewInt(3000000))

	if err != nil {
		log.Println(err)
	}
	if err == nil {
		log.Println(tx.Hash())
	}
}
func TestTetherTokenFilterer_WatchApproval(t *testing.T) {

	conn, err := ethclient.Dial("https://mainnet.infura.io/v3/dab126a4e1f444569c8f517a42cddda2")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum web3client: %v", err)
	}
	store, err := NewTetherTokenCaller(common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Storage contract: %v", err)
	}
	amount, err := store.BalanceOf(nil, common.HexToAddress("0xE3DB91F3B07F282f9890aE52e7fC0Bd2Adee2C74"))
	log.Println(amount)
	approvalAmount, err := store.Allowed(nil, common.HexToAddress("0xE3DB91F3B07F282f9890aE52e7fC0Bd2Adee2C74"), common.HexToAddress("0xFF8C78F54235D04f88378f98f5ec6Fa68802f1b4"))
	log.Println(approvalAmount)

	log.Println(approvalAmount.BitLen())
	if approvalAmount.BitLen() > 0 {
		log.Println("approval ulimit")
	} else {

		log.Println("111111111111111111")
	}
	//// Create an authorized transactor and call the store function
	//nonce, _ := conn.NonceAt(context.Background(), common.HexToAddress("你私钥对应的账户地址"), nil)
	//gasPrice, _ := conn.SuggestGasPrice(context.Background())
	////用哪条链，就用那个id
	//auth, err := bind.NewKeyedTransactorWithChainID(PrivateKey, big.NewInt(5))
	//auth.GasLimit = uint64(300000)
	//auth.Nonce = new(big.Int).SetUint64(nonce)
	//auth.GasPrice = gasPrice
	//if err != nil {
	//	log.Fatalf("Failed to create authorized transactor: %v", err)
	//}
	//// Call the store() function
	//tx, err := store.Store(auth, big.NewInt(420))
	//if err != nil {
	//	log.Fatalf("Failed to update value: %v", err)
	//}
	//fmt.Printf("Update pending: 0x%x\n", tx.Hash())

}
