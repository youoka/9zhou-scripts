package internal

import (
	"9zhou-scripts/client"
	"errors"
	"fmt"
	"math"
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
func (s *Service) Pay(num1000, num500, num200, num100 int) error {
	f, err := strconv.ParseFloat(s.shop.UserInfo.Data.Wallet.Balance, 64)
	if err != nil {
		fmt.Println("余额转换错误:", err)
		return err
	}
	if f < 1000 {
		fmt.Println("余额不足")
		return errors.New("余额不足")
	}
	fmt.Println("余额:", f)
	fmt.Println("开始下单")
	count := math.Floor(f / 1000)
	for i := 0; i < int(count); i++ {
		time.Sleep(time.Second * 2)
		orderId, err := s.shop.CreateOrder(client.P1000)
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
	}
	fmt.Println("下单完成")
	return nil
}
