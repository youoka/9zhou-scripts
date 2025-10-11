package config

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type Config struct {
	ShopAccount     string `yaml:"shop_account"`
	ShopPassword    string `yaml:"shop_password"`
	ReclaimAccount  string `yaml:"reclaim_account"`
	ReclaimPassword string `yaml:"reclaim_password"`
	AuthKey         string `yaml:"auth_key"`
}

func LoadConfig() (*Config, error) {
	// 检查key.yaml是否存在
	if _, err := os.Stat("key.yaml"); os.IsNotExist(err) {
		// 文件不存在，需要用户输入信息
		return createConfigFromInput()
	}

	// 文件存在，从文件中加载配置
	data, err := os.ReadFile("key.yaml")
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func createConfigFromInput() (*Config, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("请输入商城账号: ")
	shopAccount, _ := reader.ReadString('\n')
	shopAccount = strings.TrimSpace(shopAccount)

	fmt.Print("请输入商城密码: ")
	shopPassword, _ := reader.ReadString('\n')
	shopPassword = strings.TrimSpace(shopPassword)

	fmt.Print("请输入核销账号: ")
	reclaimAccount, _ := reader.ReadString('\n')
	reclaimAccount = strings.TrimSpace(reclaimAccount)

	fmt.Print("请输入核销密码: ")
	reclaimPassword, _ := reader.ReadString('\n')
	reclaimPassword = strings.TrimSpace(reclaimPassword)

	fmt.Print("请输入授权密钥: ")
	authKey, _ := reader.ReadString('\n')
	authKey = strings.TrimSpace(authKey)

	config := &Config{
		ShopAccount:     shopAccount,
		ShopPassword:    shopPassword,
		ReclaimAccount:  reclaimAccount,
		ReclaimPassword: reclaimPassword,
		AuthKey:         authKey,
	}

	// 保存到key.yaml
	err := saveConfig(config)
	if err != nil {
		return nil, err
	}

	fmt.Println("配置已保存到key.yaml")

	return config, nil
}

func saveConfig(config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile("key.yaml", data, 0644)
}
