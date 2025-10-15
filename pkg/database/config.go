package database

type Config struct {
	Type    int    `gorm:"column:type;primary_key" json:"type"`
	Num1000 int    `gorm:"column:num1000" json:"num1000"`
	Num500  int    `gorm:"column:num500" json:"num500"`
	Num200  int    `gorm:"column:num200" json:"num200"`
	Num100  int    `gorm:"column:num100" json:"num100"`
	Date    string `gorm:"column:date" json:"date"`
}

func (c Config) TableName() string {
	return "config"
}
func (d *DataBase) CreateConfig(c *Config) error {
	return d.conn.Create(c).Error
}
func (d *DataBase) UpdateConfig(c *Config) error {
	return d.conn.Save(c).Error
}
func (d *DataBase) GetConfig() (*Config, error) {
	var c Config
	err := d.conn.Where("type = ?", 1).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
