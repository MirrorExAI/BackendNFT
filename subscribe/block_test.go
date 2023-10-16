package subscribe

import "testing"

func TestSubscribeBlock_Subscribe(t *testing.T) {
	subBlock, _ := NewSubscribeBlock("wss://mainnet.infura.io/ws/v3/dab126a4e1f444569c8f517a42cddda2")
	subBlock.Subscribe()
}
