package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type SysCheckInUser struct {
	global.GVA_MODEL
	CheckDate string `json:"check_date" gorm:"comment:打卡日期"` //
	Phone     string `json:"phone" gorm:"comment:号码"`        //
	Email     string `json:"email" gorm:"comment:邮箱"`        //

}

func (SysCheckInUser) TableName() string {
	return "sys_check_in_users"
}
