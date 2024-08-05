package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id        int    `orm:"auto"`
	Username  string `orm:"size(100)"`
	Password  string `orm:"size(100)"`
	Email     string `orm:"size(100)"`
	Role      string `orm:"size(100)"`
	OtpSecret string `orm:"size(100)"`
}

func init() {
	// Đăng ký model
	orm.RegisterModel(new(User))
}
