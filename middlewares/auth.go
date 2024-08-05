package middlewares

import (
	"base-beego-project/utils" // Thay thế bằng tên dự án của bạn
	"fmt"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

func AuthMiddleware(ctx *context.Context) {
	authHeader := ctx.Input.Header("Authorization")
	fmt.Println("authHeader: " + authHeader)
	if authHeader == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Authorization header missing"}, false, false)
		return
	}

	token := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Invalid token"}, false, false)
		return
	}

	ctx.Input.SetData("username", claims.Username)
	ctx.Input.SetData("role", claims.Role)
}

func AdminOnly(ctx *context.Context) {
	role := ctx.Input.GetData("role")
	if role != "admin" {
		ctx.Output.SetStatus(http.StatusForbidden)
		ctx.Output.JSON(map[string]string{"error": "Access denied"}, false, false)
	}
}
