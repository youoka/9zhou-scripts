package database

type HxAccount struct {
	Type     int    `gorm:"column:type;primary_key" json:"type"`
	Account  string `gorm:"column:account;unique;not null;type:text" json:"account"`
	Password string `gorm:"column:password;type:text" json:"password"`
	Key      string `gorm:"column:key;type:text" json:"key"`
}

func (h HxAccount) TableName() string {
	return "hx_account"
}
func (d *DataBase) CreateHxAccount(account *HxAccount) error {
	account.Type = 1
	return d.conn.Create(account).Error
}

func (d *DataBase) UpdateHxAccount(account *HxAccount) error {
	account.Type = 1
	return d.conn.Model(account).Updates(account).Error
}

func (d *DataBase) GetHxAccount() (*HxAccount, error) {
	var account HxAccount
	err := d.conn.Where("type = ?", 1).First(&account).Error
	return &account, err
}
