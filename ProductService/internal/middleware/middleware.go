package middleware

import (
	"ProductService/config"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"net/http"
)

func DefMiddleware(e *echo.Echo) {

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

}

func ValidateToken(token string) bool {
	url := fmt.Sprintf("http://%s/auth/validate", config.ServiceConfig.AUTH_SERVICE_URL)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("AuthService unavailable:", err)
		return false
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("AuthService response:", string(body))

	return resp.StatusCode == http.StatusOK
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" || !ValidateToken(token) {
			return c.JSON(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}
