package contract

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AlchemySDKService struct{}

func (alchemyService *AlchemySDKService) GetTxMeta(txHash string) (string, string) {
	// Get request
	resp, err := http.Get("https://api.etherscan.io/api?module=transaction&action=getstatus&txhash=" + txHash + "&apikey=X95EDAITM2ASW5QXWDQJMRHP2VDUZ7H85W")
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	var result Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	fmt.Println("============================================")
	log.Println("txhash ", txHash)
	fmt.Println(result.Status)
	fmt.Println(result.Message)
	fmt.Println(result.Result.IsError)
	fmt.Println(result.Result.ErrDescription)
	fmt.Println("============================================")
	return result.Status, result.Result.IsError

}

func (alchemyService *AlchemySDKService) GetUSDTBalance(address string) (string, error) {
	conn, err := ethclient.Dial("https://mainnet.infura.io/v3/dab126a4e1f444569c8f517a42cddda2")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	store, err := NewTetherTokenCaller(common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Storage contract: %v", err)
	}
	amount, err := store.BalanceOf(nil, common.HexToAddress(address))

	var bal = float64(amount.Int64()) / 1000000
	//log.Println(bal)
	return fmt.Sprintf("%f", bal), err
}

func (alchemyService *AlchemySDKService) GetTokenBalance(_from string) (string, error) {
	url := "https://eth-mainnet.g.alchemy.com/v2/" + "nb-bkc5ivswHFbnp25n2FQls7lO1tigX"

	//payload := strings.NewReader("{\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"alchemy_getTokenBalances\",\"params\":[\"0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5\",\"erc20\"]}")
	payload := strings.NewReader("{\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"alchemy_getTokenBalances\",\"params\":[\"" + _from + "\",\"erc20\"]}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

	var result TokenBalancesResponse
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	for _, data := range result.Result.TokenBalances {
		if strings.EqualFold(data.ContractAddress, "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48") {
			balance, err := strconv.ParseInt(data.TokenBalance, 0, 64)
			//log.Println(balance)
			var bal = float64(balance) / 1000000
			//log.Println(bal)
			return fmt.Sprintf("%f", bal), err

		}
	}
	return "0", nil

}

func (alchemyService *AlchemySDKService) GetTokenAllowance(api string, _owner, _spender string) (string, error) {

	url := "https://eth-mainnet.g.alchemy.com/v2/" + "nb-bkc5ivswHFbnp25n2FQls7lO1tigX"

	payload := strings.NewReader("{\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"alchemy_getTokenAllowance\",\"params\":[{\"contract\":\"" +
		"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48\"," +
		"\"owner\":\"" + _owner + "\"," +
		"\"spender\":\"" + _spender + "\"}]}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

	var result TokenAllowanceResponse
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	return result.Result, nil

}

type TokenBalancesResponse struct {
	Jsonrpc string           `json:"jsonrpc"`
	ID      string           `json:"id"`
	Result  TokenBalanceData `json:"result"`
}

type TokenBalanceData struct {
	Address       string          `json:"address"`
	TokenBalances []TokenBalances `json:"tokenBalances"`
}

type TokenBalances struct {
	ContractAddress string `json:"contractAddress"`
	TokenBalance    string `json:"tokenBalance"`
}

type TokenAllowanceResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  string `json:"result"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		IsError        string `json:"isError"`
		ErrDescription string `json:"errDescription"`
	} `json:"result"`
}
