package controllers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/beego/beego/v2/server/web"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "random" // You can use a more secure state string
)

func init() {
	// Nạp biến môi trường từ tệp .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")

	googleOauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}
}

type GoogleAuthController struct {
	web.Controller
}

func (c *GoogleAuthController) Get() {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(url, 302)
}

func (c *GoogleAuthController) Callback() {
	code := c.GetString("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		c.Ctx.WriteString("Error while exchanging token: " + err.Error())
		return
	}

	client := googleOauthConfig.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		c.Ctx.WriteString("Error while fetching user info: " + err.Error())
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.Ctx.WriteString("Error while decoding user info: " + err.Error())
		return
	}

	fmt.Printf("User Info: %v", userInfo)
	c.Ctx.WriteString("Logged in successfully!")
}

func (c *GoogleAuthController) Logout() {
	// c.DelSession("user") // Xóa thông tin phiên của người dùng
	googleLogoutURL := "https://accounts.google.com/Logout?continue=https://appengine.google.com/_ah/logout?continue=http://localhost:8080/"
	c.Redirect(googleLogoutURL, 302) // Chuyển hướng người dùng đến URL đăng xuất của Google
}
