package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type VictimRouter struct{}

func (s *VictimRouter) InitVictimRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {
	VictimRouter := Router.Group("victim").Use(middleware.OperationRecord())
	//VictimRouter := Router.Group("victim")
	victimRouter := v1.ApiGroupApp.SystemApiGroup.SystemVictim
	victimTxRouter := v1.ApiGroupApp.SystemApiGroup.SystemVictimTx
	{
		VictimRouter.POST("refreshTx", victimTxRouter.Refresh)           // 创建Victim
		VictimRouter.POST("refreshBalance", victimRouter.RefreshBalance) // 创建Victim

		VictimRouter.POST("createVictim", victimRouter.CreateVictim)               // 创建Victim
		VictimRouter.POST("createVictimTx", victimTxRouter.CreateVictimTx)         // 创建VictimTx
		VictimRouter.POST("getVictimTxList", victimTxRouter.GetVictimTxList)       // 创建VictimTx
		VictimRouter.POST("getSumVictimTxList", victimTxRouter.GetSumVictimTxList) // 创建VictimTx
		VictimRouter.POST("getVictimTxById", victimTxRouter.GetVictimTxById)       // 获取单条Victim消息
		VictimRouter.POST("statVictimTx", victimTxRouter.StatVictimTxInfo)         // 创建VictimTx
		VictimRouter.POST("statVictimTx2", victimTxRouter.StatVictimTx2)           // 创建VictimTx
		VictimRouter.POST("deleteVictim", victimRouter.DeleteVictim)               // 删除Victim
		VictimRouter.POST("getVictimById", victimRouter.GetVictimById)             // 获取单条Victim消息

		VictimRouter.POST("updateVictim", victimRouter.UpdateVictim)               // 更新Victim
		VictimRouter.DELETE("deleteVictimsByIds", victimRouter.DeleteVictimsByIds) // 删除选中Victim
		VictimRouter.POST("getAllVictims", victimRouter.GetAllVictims)             // 获取所有Victim
		VictimRouter.POST("getVictimList", victimRouter.GetVictimList)             // 获取Victim列表

		VictimRouter.POST("getUsers", victimRouter.GetUsers) // 获取自身信息
	
	}

}
