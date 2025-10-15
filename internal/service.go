package internal

import (
	"9zhou-scripts/client"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Service struct {
	shop    *client.ShopAccount
	reclaim *client.ReclaimAccount
}

func NewService(shop *client.ShopAccount, reclaim *client.ReclaimAccount) *Service {
	return &Service{shop: shop, reclaim: reclaim}
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
	err := s.shop.Login()
	if err != nil {
		return err
	}
	err = s.shop.Info()
	if err != nil {
		return err
	}
	err = s.reclaim.Login()
	if err != nil {
		return err
	}
	err = s.reclaim.Info()
	if err != nil {
		return err
	}
	return err
}
func (s *Service) Transfer() error {
	return s.reclaim.Transfer(s.shop.Account)
}
func (s *Service) Pay(p int, num int) error {
	s.shop.Info()
	balance, err := strconv.ParseFloat(s.shop.UserInfo.Data.Wallet.Balance, 64)
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
		time.Sleep(time.Second * 2)
		orderId, err := s.shop.CreateOrder(goods)
		if err != nil {
			return err
		}
		if orderId != "" {
			msg, err := s.shop.PayOrder(orderId)
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
//	f, err := strconv.ParseFloat(s.shop.UserInfo.Data.Wallet.Balance, 64)
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
//		orderId, err := s.shop.CreateOrder(client.P1000)
//		if err != nil {
//			return err
//		}
//		if orderId != "" {
//			msg, err := s.shop.PayOrder(orderId)
//			fmt.Println(msg)
//			if err != nil {
//				return err
//			}
//		}
//	}
//	fmt.Println("下单完成")
//	return nil
//}
