package db

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

type TinyUrl struct {
	gorm.Model
	LongUrl    string `gorm:"size:200"`
	ShortUrl   string `gorm:"size:50"`
	ExpireTime time.Time
	Counter    float64
}

type URLDetail struct {
	URL                 string        `json:"url"`
	CreatedAt           string        `json:"created_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

func InitDb() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "tinyurl:tinyurl@tcp(42.192.36.14:3306)/tinyurl?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := DB.AutoMigrate(&TinyUrl{}); err != nil {
		panic(err)
	}
}
