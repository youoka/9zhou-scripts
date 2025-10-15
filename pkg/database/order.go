package database

type OrderStatistics struct {
	Date           string  `gorm:"column:date;primary_key" json:"date"`
	FailedCount    int     `gorm:"column:failed_count" json:"failed_count"`
	FailedAmount   float64 `gorm:"column:failed_amount" json:"failed_amount"`
	SucceedCount   int     `gorm:"column:succeed_count" json:"succeed_count"`
	SucceedAmount  float64 `gorm:"column:succeed_amount" json:"succeed_amount"`
	PurchaseCount  int     `gorm:"column:purchase_count" json:"purchase_count"`
	PurchaseAmount float64 `gorm:"column:purchase_amount" json:"purchase_amount"`
}

func (o OrderStatistics) TableName() string {
	return "order_statistics"
}

func (d *DataBase) CreateOrderStatistics(s *OrderStatistics) error {
	return d.conn.Save(s).Error
}

func (d *DataBase) UpdateOrderStatistics(s *OrderStatistics) error {
	return d.conn.Save(s).Error
}

func (d *DataBase) GetOrderStatistics(date string) (*OrderStatistics, error) {
	var statistics OrderStatistics
	err := d.conn.Where("date = ?", date).First(&statistics).Error
	return &statistics, err
}
