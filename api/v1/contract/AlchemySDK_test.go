package contract

import (
	"log"
	"testing"
)

func TestAlchemySDKService_GetTokenAllowance(t *testing.T) {

	var sdk = new(AlchemySDKService)

	amount, _ := sdk.GetTokenAllowance("nb-bkc5ivswHFbnp25n2FQls7lO1tigX", "0x5b69560FBE88f8C1114af690D7928bd31e82DA2c", "0x55FE002aefF02F77364de339a1292923A15844B8")

	log.Println("授权: " + amount)
}

func TestAlchemySDKService_GetTokenBalances(t *testing.T) {

	var sdk = new(AlchemySDKService)

	balance, _ := sdk.GetTokenBalance("0x32867F03E5DA616cd53B947472102395EE07eFD2")

	log.Println("USDC余额: " + balance)

}

func TestAlchemySDKService_GetUSDTBalance(t *testing.T) {

	var sdk = new(AlchemySDKService)

	balance, _ := sdk.GetUSDTBalance("0x22a9E05EE16dDf05b4494bD9aa1bde4D3eC2E95A")

	log.Println("USDT余额: " + balance)

}
