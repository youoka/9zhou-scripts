package service

import (
	"9zhou-scripts/client"
	"9zhou-scripts/pkg/database"
	"9zhou-scripts/pkg/utils"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var ZhouService *Service

func Init() {
	hxAccount, err := database.Db.GetHxAccount()
	if err != nil {
		fmt.Println("无法获取核销账号信息")
		return
	}

	// 获取商店账号信息（获取第一个商店账号）
	shopAccounts, err := database.Db.GetAllShopAccount()
	if err != nil || len(shopAccounts) == 0 {
		fmt.Println("无法获取商城账号信息")
		return
	}
	shopAccount := shopAccounts[0]
	// 授权码验证
	if !utils.CheckIsValid(hxAccount.Key, shopAccount.Account, hxAccount.Account) {
		fmt.Println("授权码无效")
		return
	}
	shop := client.NewShopAccount(shopAccount.Account, shopAccount.Password)
	reclaim := client.NewReclaimAccount(hxAccount.Account, hxAccount.Password)
	ZhouService = NewService(shop, reclaim)
}

type Service struct {
	Shop    *client.ShopAccount
	Reclaim *client.ReclaimAccount
}

func NewService(shop *client.ShopAccount, reclaim *client.ReclaimAccount) *Service {
	return &Service{Shop: shop, Reclaim: reclaim}
}

func (s *Service) StartAuto2Task() {
	go func() {
		ticker := time.NewTicker(3 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.Login()
			}
		}
	}()
}
func (s *Service) Login() error {
	err := s.Shop.Login()
	if err != nil {
		return err
	}
	err = s.Shop.Info()
	if err != nil {
		return err
	}
	err = s.Reclaim.Login()
	if err != nil {
		return err
	}
	err = s.Reclaim.Info()
	if err != nil {
		return err
	}
	return err
}
func (s *Service) Transfer() error {
	return s.Reclaim.Transfer(s.Shop.Account)
}
func (s *Service) Pay(p int, num int) error {
	s.Shop.Info()
	balance, err := strconv.ParseFloat(s.Shop.UserInfo.Data.Wallet.Balance, 64)
	if err != nil {
		fmt.Println("余额转换错误:", err)
		return err
	}
	if balance < float64(p*num) {
		fmt.Println("余额不足")
		return errors.New("余额不足")
	}
	goods := ""
	switch p {
	case 1000:
		goods = client.P1000
	case 500:
		goods = client.P500
	case 200:
		goods = client.P200
	case 100:
		goods = client.P100
	default:
		return errors.New("商品不存在")
	}
	for i := range num {
		time.Sleep(time.Second * 1)
		orderId, err := s.Shop.CreateOrder(goods)
		if err != nil {
			return err
		}
		if orderId != "" {
			msg, err := s.Shop.PayOrder(orderId)
			fmt.Println(msg)
			if err != nil {
				return err
			}
		}
		fmt.Println(fmt.Sprintf("第%d个商品下单成功", i+1))
	}
	return nil
}
func (s *Service) StartPay(num1000, num500, num200, num100 int) {
	s.Pay(1000, num1000)
	s.Pay(500, num500)
	s.Pay(200, num200)
	s.Pay(100, num100)
}

//func (s *Service) Pay(num1000, num500, num200, num100 int) error {
//	f, err := strconv.ParseFloat(s.Shop.UserInfo.Data.Wallet.Balance, 64)
//	if err != nil {
//		fmt.Println("余额转换错误:", err)
//		return err
//	}
//	if f < 1000 {
//		fmt.Println("余额不足")
//		return errors.New("余额不足")
//	}
//	fmt.Println("余额:", f)
//	fmt.Println("开始下单")
//	count := math.Floor(f / 1000)
//	for i := 0; i < int(count); i++ {
//		time.Sleep(time.Second * 2)
//		orderId, err := s.Shop.CreateOrder(client.P1000)
//		if err != nil {
//			return err
//		}
//		if orderId != "" {
//			msg, err := s.Shop.PayOrder(orderId)
//			fmt.Println(msg)
//			if err != nil {
//				return err
//			}
//		}
//	}
//	fmt.Println("下单完成")
//	return nil
//}
