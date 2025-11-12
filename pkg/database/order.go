package database

type OrderStatistics struct {
	Date        string  `gorm:"column:date;primary_key" json:"date"`
	Account     string  `gorm:"column:account;primary_key" json:"account"`
	Num1000     int     `gorm:"column:num1000" json:"num1000"`
	Num500      int     `gorm:"column:num500" json:"num500"`
	Num200      int     `gorm:"column:num200" json:"num200"`
	Num100      int     `gorm:"column:num100" json:"num100"`
	TotalAmount float64 `gorm:"column:total_amount" json:"total_amount"`
}

func (o OrderStatistics) TableName() string {
	return "order_statistics"
}

func (d *DataBase) SaveOrderStatistics(s *OrderStatistics) error {
	return d.conn.Save(s).Error
}
func (d *DataBase) GetOrderStatistics(date string) (*OrderStatistics, error) {
	var statistics OrderStatistics
	err := d.conn.Where("date = ?", date).First(&statistics).Error
	return &statistics, err
}
