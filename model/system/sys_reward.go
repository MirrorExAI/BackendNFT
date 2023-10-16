package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type SysReward struct {
	global.GVA_MODEL

	Reward  string `json:"reward" gorm:"comment:奖励"`  // 奖励
	Address string `json:"address" gorm:"comment:地址"` // 地址
	Agent   string `json:"agent" gorm:"comment:渠道"`   // 渠道代理
	Status  uint   `json:"status" gorm:"comment:状态"`  // 状态

}

func (SysReward) TableName() string {
	return "sys_reward"
}
