package main

import (
	_ "base-beego-project/models"  // Thay thế "base-beego-project" bằng tên dự án của bạn
	_ "base-beego-project/routers" // Thay thế "base-beego-project" bằng tên dự án của bạn
	_ "base-beego-project/utils"   // Thay thế "base-beego-project" bằng tên dự án của bạn
	"log"
	"os"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	// Nạp biến môi trường từ tệp .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Lấy biến môi trường
	driver := os.Getenv("DB_DRIVER")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	// Tạo chuỗi kết nối
	connStr := "user=" + user + " password=" + password + " dbname=" + dbname + " host=" + host + " port=" + port + " sslmode=" + sslmode

	// Đăng ký cơ sở dữ liệu PostgreSQL
	orm.RegisterDriver(driver, orm.DRPostgres)
	orm.RegisterDataBase("default", driver, connStr)

	// Tạo bảng tự động
	orm.RunSyncdb("default", false, true)
}

func main() {
	web.Run()
}
