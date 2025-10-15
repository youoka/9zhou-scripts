package database

type ShopAccount struct {
	Account  string `gorm:"column:account;unique;not null;type:text" json:"account"`
	Password string `gorm:"column:password;type:text" json:"password"`
}

func (s ShopAccount) TableName() string {
	return "shop_account"
}
func (d *DataBase) CreateShopAccount(account *ShopAccount) error {
	return d.conn.Create(account).Error
}
func (d *DataBase) GetAllShopAccount() ([]*ShopAccount, error) {
	var shopAccounts []*ShopAccount
	err := d.conn.Find(&shopAccounts).Error
	if err != nil {
		return nil, err
	}
	return shopAccounts, nil
}
