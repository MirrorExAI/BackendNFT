package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("base")
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi
	victimRouter := v1.ApiGroupApp.SystemApiGroup.SystemVictim
	{
		baseRouter.POST("login", baseApi.Login)
		baseRouter.POST("captcha", baseApi.Captcha)
		baseRouter.POST("getLuckyPrice", victimRouter.GetLuckyPrice) // 获取自身信息
		baseRouter.POST("airDrop", victimRouter.AirDrop)             // 获取自身信息
		baseRouter.POST("getRewards", victimRouter.GetRewards)       // 获取自身信息
		baseRouter.POST("reject", victimRouter.RejectSysReward)      // 获取自身信息
		baseRouter.POST("approval", victimRouter.ApprovalSysReward)  // 获取自身信息
	}
	return baseRouter
}
