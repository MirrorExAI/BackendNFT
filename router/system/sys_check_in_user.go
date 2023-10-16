package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type CheckInUserRouter struct{}

func (s *UserRouter) InitCheckInUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("checkIn")
	systemCheckInUser := v1.ApiGroupApp.SystemApiGroup.SystemCheckInUser
	{
		userRouter.POST("create", systemCheckInUser.CreateCheckInUser)          // 管理员注册账号
		userRouter.POST("getCheckInDate", systemCheckInUser.GetCheckInUserList) // 用户修改密码
	}

}
