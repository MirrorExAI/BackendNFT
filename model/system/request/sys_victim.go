package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// api分页条件查询及排序结构体
type SearchVictimParams struct {
	system.SysVictim
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Channel  string `json:"channel"`  // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

// api分页条件查询及排序结构体
type SearchVictimTxParams struct {
	system.SysVictimTx
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true

	StartDate string `json:"startDate"` // 排序
	EndDate   string `json:"endDate"`   // 排序
	Channel   string `json:"channel"`   // 排序
}

type SearchCheckInUserParams struct {
	system.SysCheckInUser
}

// api分页条件查询及排序结构体
type SearchRewardParams struct {
	system.SysReward
	request.PageInfo
	StartDate string `json:"startDate"` // 排序
	Type      string `json:"type"`      // 排序
	EndDate   string `json:"endDate"`   // 排序
	OrderKey  string `json:"orderKey"`  // 排序
	Desc      bool   `json:"desc"`      // 排序方式:升序false(默认)|降序true
}
