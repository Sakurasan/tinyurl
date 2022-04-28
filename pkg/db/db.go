package db

import (
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

type TinyUrl struct {
	gorm.Model
	LongUrl    string `gorm:"size:7713"`
	ShortUrl   string `gorm:"size:20"`
	ExpireTime time.Time
	Counter    float64
	AddIP      string
}

type URLDetail struct {
	URL                 string        `json:"url"`
	CreatedAt           string        `json:"created_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

func InitDb() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	var dsn string
	if len(os.Getenv("DSN")) > 1 {
		dsn = os.Getenv("DSN")
	} else {
		dsn = "tinyurl:tinyurl@tcp(42.192.36.14:3306)/tinyurl?charset=utf8mb4&parseTime=True&loc=Local"
	}
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := DB.AutoMigrate(&TinyUrl{}); err != nil {
		panic(err)
	}
}
