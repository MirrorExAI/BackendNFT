package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

type CheckInUserService struct{}

func (CheckInUserService *CheckInUserService) CreateCheckInUser(user system.SysCheckInUser) (err error) {
	return global.GVA_DB.Create(&user).Error
}

func (CheckInUserService *CheckInUserService) GetCheckInUserList(user system.SysCheckInUser) (list interface{}, err error) {
	db := global.GVA_DB.Model(&system.SysCheckInUser{})
	var victimList []system.SysCheckInUser

	if user.Phone == "" && user.Email == "" {
		return nil, nil
	}
	if user.Phone != "" {
		db = db.Where("phone = ?", user.Phone)
	}

	if user.Email != "" {
		db = db.Where("email = ?", user.Email)
	}

	err = db.Find(&victimList).Error
	return victimList, err
}
