package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/system"

type SysVictimResponse struct {
	Victim system.SysVictim `json:"victim"`
}

type SysVictimTxResponse struct {
	Victim system.SysVictimTx `json:"victim"`
}
