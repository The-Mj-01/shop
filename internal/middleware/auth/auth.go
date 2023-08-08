package auth

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"shop/pkg/IP"
	"shop/pkg/tokenizer"
	"strings"
)

// unAuthorizedMsg defines a message for authorization
var unAuthorizedMsg map[string]string = map[string]string{
	"data":    "",
	"message": "You are not authorized",
}

// ValidateJWT for user
func ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bearerToken := c.Request().Header.Get("Authorization")
		if bearerToken == "" {
			return c.JSON(http.StatusUnauthorized, unAuthorizedMsg)
		}
		bearerToken = fixSentTokenString(bearerToken)

		ip := IP.ExtractFromEcho(c)
		isAllowed := handleJwtValidation(bearerToken, ip)

		if !isAllowed {
			return c.JSON(http.StatusUnauthorized, unAuthorizedMsg)
		}

		return next(c)
	}
}

// handleJwtValidation handles JWT verification functionality
func handleJwtValidation(token, ip string) bool {
	ctx := context.Background()
	tokenHandler := tokenizer.CreateTokenizer(ctx)
	return tokenHandler.IsActive(token, ip)
}

// fixSentTokenString and return it in order to make it use-able
func fixSentTokenString(token string) string {
	token = strings.Replace(token, "Bearer", "", 1)
	token = strings.Replace(token, "bearer", "", 1)
	token = strings.Trim(token, " ")
	return token
}
