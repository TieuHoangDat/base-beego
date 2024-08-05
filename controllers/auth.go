package controllers

import (
	"base-beego-project/models"
	"base-beego-project/utils"
	"fmt"

	// "encoding/json"
	"net/http"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	web.Controller
}

func (c *AuthController) CreateUser() {
	// Lấy dữ liệu từ request
	username := c.GetString("Username")
	password := c.GetString("Password")
	email := c.GetString("Email")
	role := c.GetString("Role")

	// Mã hóa mật khẩu
	hashedPassword, errHashedPassword := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errHashedPassword != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.ResponseWriter.Write([]byte("Error encrypting password"))
		return
	}

	// Tạo mã OTP
	otp, secret, err := utils.GenerateOTP(email)
	fmt.Println(otp) // xem otp
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.ResponseWriter.Write([]byte("Error generating OTP"))
		return
	}

	// Gửi OTP qua email
	err = utils.SendOTP(email, otp)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.ResponseWriter.Write([]byte("Error sending OTP"))
		return
	}

	// Tạo đối tượng User
	user := models.User{
		Username:  username,
		Password:  string(hashedPassword),
		Email:     email,
		Role:      role,
		OtpSecret: secret, // Lưu trữ secret để xác thực OTP
	}

	// Tạo đối tượng ORM
	o := orm.NewOrm()

	// Thêm người dùng vào cơ sở dữ liệu
	_, err = o.Insert(&user)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.ResponseWriter.Write([]byte("Error adding user"))
		return
	}

	// Phản hồi thành công dưới dạng JSON
	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
	c.Data["json"] = map[string]string{
		"message":  "User added successfully",
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	}
	c.ServeJSON()
}

func (c *AuthController) VerifyOTP() {
	email := c.GetString("Email")
	otp := c.GetString("OTP")

	// Tạo đối tượng ORM
	o := orm.NewOrm()
	user := models.User{Email: email}

	// Tìm người dùng trong cơ sở dữ liệu
	err := o.Read(&user, "Email")
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{"error": "Invalid email"}
		c.ServeJSON()
		return
	} else if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Database error"}
		c.ServeJSON()
		return
	}

	fmt.Println(email + " " + otp + " " + user.OtpSecret)

	// Xác thực OTP
	valid := totp.Validate(otp, user.OtpSecret)
	if !valid {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{"error": "Invalid OTP"}
		c.ServeJSON()
		return
	}

	// Cập nhật trạng thái xác thực thành công
	user.OtpSecret = "verified"
	_, err = o.Update(&user, "OtpSecret")
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to update OTP status"}
		c.ServeJSON()
		return
	}

	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
	c.Data["json"] = map[string]string{
		"message": "OTP verified successfully!",
	}
	c.ServeJSON()
}

func (c *AuthController) Login() {
	// Lấy dữ liệu từ request
	username := c.GetString("Username")
	password := c.GetString("Password")

	// Kiểm tra xem các trường có được cung cấp không
	if username == "" || password == "" {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Username and password required"}
		c.ServeJSON()
		return
	}

	// Tạo đối tượng ORM
	o := orm.NewOrm()
	user := models.User{Username: username}

	// Tìm người dùng trong cơ sở dữ liệu
	err := o.Read(&user, "Username")
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{"error": "Invalid username or password"}
		c.ServeJSON()
		return
	} else if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Database error"}
		c.ServeJSON()
		return
	}

	// Kiểm tra mật khẩu
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{"error": "Invalid username or password"}
		c.ServeJSON()
		return
	}

	// Kiểm tra xem OtpSecret đã được xác thực chưa
	if user.OtpSecret != "verified" {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{"error": "OTP not verified"}
		c.ServeJSON()
		return
	}

	// Sinh JWT token
	token, err := utils.GenerateJWT(user.Username, user.Role)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Error generating token"}
		c.ServeJSON()
		return
	}

	// Sinh Refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.Username, user.Role)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Error generating refresh token"}
		c.ServeJSON()
		return
	}

	// Phản hồi thành công dưới dạng JSON
	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
	c.Data["json"] = map[string]interface{}{
		"message":      "Login successful",
		"username":     user.Username,
		"email":        user.Email,
		"role":         user.Role,
		"token":        token,
		"refreshToken": refreshToken,
	}
	c.ServeJSON()
}

func (c *AuthController) RefreshToken() {
	refreshToken := c.GetString("refreshToken")

	claims, err := utils.ValidateJWT(refreshToken)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{"error": "Invalid refresh token"}
		c.ServeJSON()
		return
	}

	token, err := utils.GenerateJWT(claims.Username, claims.Role)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Error generating new token"}
		c.ServeJSON()
		return
	}

	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
	c.Data["json"] = map[string]string{
		"token": token,
	}
	c.ServeJSON()
}
