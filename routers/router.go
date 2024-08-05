package routers

import (
	"base-beego-project/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// auth
	beego.Router("/user", &controllers.AuthController{}, "post:CreateUser")
	beego.Router("/login", &controllers.AuthController{}, "post:Login")
	beego.Router("/user/refresh-token", &controllers.AuthController{}, "post:RefreshToken")
	beego.Router("/verifyotp", &controllers.AuthController{}, "post:VerifyOTP")

	// google auth
	beego.Router("/auth", &controllers.GoogleAuthController{})
	beego.Router("/auth/callback", &controllers.GoogleAuthController{}, "get:Callback")
	beego.Router("/auth/logout", &controllers.GoogleAuthController{}, "get:Logout") // http://localhost:8080/auth/logout

	beego.Router("/user", &controllers.UserController{}, "get:GetAllUsers")
}
