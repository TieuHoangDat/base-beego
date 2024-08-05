package controllers

import (
	"base-beego-project/middlewares" // Thêm import cho middlewares
	"base-beego-project/models"

	// "encoding/json"
	"net/http"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

type UserController struct {
	web.Controller
}

func (c *UserController) Prepare() {
	middlewares.AuthMiddleware(c.Ctx)
}

func (c *UserController) GetAllUsers() {
	// Check for admin role
	middlewares.AdminOnly(c.Ctx) // Kiểm tra quyền admin

	o := orm.NewOrm()
	var users []models.User
	_, err := o.QueryTable(new(models.User)).All(&users)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to retrieve users"}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = users
	c.ServeJSON()
}
