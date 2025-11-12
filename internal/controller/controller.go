package controller

import (
	"9zhou-scripts/internal/service"
	"9zhou-scripts/pkg/database"
	"9zhou-scripts/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func Init(ctx *gin.Context) {
	service.Init()
	ctx.JSON(200, gin.H{
		"message": "初始化成功",
	})
}
func Login(ctx *gin.Context) {
	err := service.ZhouService.Login()
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	info := service.ZhouService.Shop.UserInfo
	reclaimInfo := service.ZhouService.Reclaim.ReclaimInfo
	ctx.JSON(200, gin.H{
		"message": "登录成功",
		"data": gin.H{
			"商城账号": info.Data.Account,
			"商城余额": info.Data.Wallet.Balance,
			"核销账号": reclaimInfo.Data.Account,
			"核销余额": reclaimInfo.Data.Wallet.Balance,
		},
	})
}
func Hx(ctx *gin.Context) {
	// 获取配置信息
	config, err := database.Db.GetConfig()
	if err != nil {
		fmt.Println("无法获取配置信息:", err)
		ctx.JSON(500, gin.H{
			"message": "无法获取配置信息",
		})
		return
	}

	// 解析配置中的日期作为开始时间
	startTime, err := time.Parse("2006-01-02", config.Date)
	if err != nil {
		fmt.Println("时间格式错误:", err)
		ctx.JSON(500, gin.H{
			"message": "时间格式错误",
		})
		return
	}

	// 结束时间是开始时间的第二天
	endTime := startTime.AddDate(0, 0, 1)
	cancelledAmount := 0.00
	cancelledCount := 0
	c := 1
	for {
		cancelledOrder, err := service.ZhouService.Shop.GetCancelledOrder(startTime.Format("2006-01-02"), endTime.Format("2006-01-02"), c)
		if err != nil {
			fmt.Println("获取取消订单失败:", err.Error())
			ctx.JSON(500, gin.H{
				"message": "获取取消订单失败",
			})
			return
		}
		cancelledAmount += utils.SumAmount(cancelledOrder)
		cancelledCount += len(cancelledOrder.Data)
		if len(cancelledOrder.Data) < 100 {
			break
		}
	}
	hxOrder := make([]string, 0)
	ShippedAmount := 0.00
	ShippedCount := 0
	i := 1
	for {
		shippedOrder, err := service.ZhouService.Shop.GetShippedOrder(startTime.Format("2006-01-02"), endTime.Format("2006-01-02"), i)
		if err != nil {
			fmt.Println("获取发货订单失败:", err.Error())
			ctx.JSON(500, gin.H{
				"message": "获取发货订单失败",
			})
			return
		}
		ShippedAmount += utils.SumAmount(shippedOrder)
		ShippedCount += len(shippedOrder.Data)
		for _, v := range shippedOrder.Data {
			hxOrder = append(hxOrder, v.Id)
		}
		i++
		if len(hxOrder) < 100 {
			break
		}
	}
	err = service.ZhouService.Reclaim.Hx(hxOrder)
	if err != nil {
		fmt.Println("核销失败:", err.Error())
		ctx.JSON(500, gin.H{
			"message": "核销失败",
		})
		return
	}
	database.Db.SaveOrderStatistics(&database.OrderStatistics{
		Date:          time.Now().Format("2006-01-02"),
		FailedCount:   cancelledCount,
		FailedAmount:  cancelledAmount,
		SucceedCount:  ShippedCount,
		SucceedAmount: ShippedAmount,
	})
	ctx.JSON(200, gin.H{
		"message": "成功",
		"data": gin.H{
			"成功订单数": ShippedCount,
			"成功金额":  ShippedAmount,
			"取消订单数": cancelledCount,
			"取消金额":  cancelledAmount,
		},
	})
}

func Transfer(ctx *gin.Context) {
	err := service.ZhouService.Transfer()
	if err != nil {
		fmt.Println("转账失败:", err.Error())
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "成功",
	})
}

func Pay(ctx *gin.Context) {
	config, err := database.Db.GetConfig()
	if err != nil {
		fmt.Println("无法获取配置信息:", err)
		ctx.JSON(500, gin.H{
			"message": "无法获取配置信息",
		})
		return
	}
	service.ZhouService.StartPay(config.Num1000, config.Num500, config.Num200, config.Num100)
	ctx.JSON(200, gin.H{
		"message": "成功",
	})
}
