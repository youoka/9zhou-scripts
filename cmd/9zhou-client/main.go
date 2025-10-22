package main

import (
	"9zhou-scripts/client"
	"9zhou-scripts/internal/service"
	"9zhou-scripts/pkg/database"
	"9zhou-scripts/pkg/utils"
	"fmt"
	"github.com/eiannone/keyboard"
	"time"
)

func main() {
	// 从数据库读取账户信息
	// 获取核销账号信息
	hxAccount, err := database.Db.GetHxAccount()
	if err != nil {
		fmt.Println("无法获取核销账号信息:", err)
		return
	}

	// 获取商店账号信息（获取第一个商店账号）
	shopAccounts, err := database.Db.GetAllShopAccount()
	if err != nil || len(shopAccounts) == 0 {
		fmt.Println("无法获取商店账号信息:", err)
		return
	}

	// 获取配置信息
	config, err := database.Db.GetConfig()
	if err != nil {
		fmt.Println("无法获取配置信息:", err)
		return
	}

	// 解析配置中的日期作为开始时间
	startTime, err := time.Parse("2006-01-02", config.Date)
	if err != nil {
		fmt.Println("时间格式错误:", err)
		return
	}

	// 结束时间是开始时间的第二天
	endTime := startTime.AddDate(0, 0, 1)
	shopAccount := shopAccounts[0]

	// 授权码验证
	if !utils.CheckIsValid(hxAccount.Key, shopAccount.Account, hxAccount.Account) {
		fmt.Println("授权码无效")
		return
	}
	fmt.Println("授权码有效")

	shop := client.NewShopAccount(shopAccount.Account, shopAccount.Password)
	reclaim := client.NewReclaimAccount(hxAccount.Account, hxAccount.Password)
	service := service.NewService(shop, reclaim)
	err = service.Login()
	if err != nil {
		fmt.Println("登录失败:", err.Error())
		return
	}
	go func() {
		service.StartAuto2Task()
	}()
	cancelledAmount := 0.00
	cancelledCount := 0
	c := 1
	for {
		cancelledOrder, err := shop.GetCancelledOrder(startTime.Format("2006-01-02"), endTime.Format("2006-01-02"), c)
		if err != nil {
			fmt.Println("获取取消订单失败:", err.Error())
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
		shippedOrder, err := shop.GetShippedOrder(startTime.Format("2006-01-02"), endTime.Format("2006-01-02"), i)
		if err != nil {
			fmt.Println("获取发货订单失败:", err.Error())
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
	err = reclaim.Hx(hxOrder)
	if err != nil {
		fmt.Println("核销失败:", err.Error())
		return
	}
	//service.Hx(startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))
	err = service.Transfer()
	if err != nil {
		fmt.Println("转账失败:", err.Error())
		return
	}
	service.StartPay(config.Num1000, config.Num500, config.Num200, config.Num100)
	PurchaseAmount := 0.00
	PurchaseCount := 0
	p := 1
	for {
		paidOrder, err := shop.GetPaidOrder(time.Now().Format("2006-01-02"), time.Now().Add(time.Hour*24).Format("2006-01-02"), p)
		if err != nil {
			fmt.Println("获取支付订单失败:", err.Error())
			return
		}
		PurchaseAmount += utils.SumAmount(paidOrder)
		p++
		PurchaseCount += len(paidOrder.Data)
		if len(paidOrder.Data) < 100 {
			break
		}
	}
	database.Db.CreateOrderStatistics(&database.OrderStatistics{
		Date:           time.Now().Format("2006-01-02"),
		FailedCount:    cancelledCount,
		FailedAmount:   cancelledAmount,
		SucceedCount:   ShippedCount,
		SucceedAmount:  ShippedAmount,
		PurchaseCount:  PurchaseCount,
		PurchaseAmount: PurchaseAmount,
	})
	// 添加任意键退出功能
	fmt.Println("按任意键退出...")
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	_, _, _ = keyboard.GetKey()
	fmt.Println("程序已退出")
}
