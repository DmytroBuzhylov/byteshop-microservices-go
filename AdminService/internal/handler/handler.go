package handler

import (
	"AdminService/config"
	"ByteShop/generated/auth"
	"ByteShop/generated/product"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

type ChangeRoleRequest struct {
	Role string `json:"role"`
}

type BanRequest struct {
	Banned bool `json:"banned"`
}

type AdminHandler struct {
	GrpcClientProductService product.ProductServiceClient
	GrpcClientAuthService    auth.AuthServiceClient
}

func (a *AdminHandler) RequestDeleteProductHandler(c echo.Context) error {

	productID := c.Param("id")

	resp, err := a.GrpcClientProductService.DeleteProduct(context.Background(), &product.DeleteProductRequest{
		ProductId: productID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": resp.Status})
}

func (a *AdminHandler) RequestDeleteAllProductsHandler(c echo.Context) error {
	userId := c.Param("id")
	resp, err := a.GrpcClientProductService.DeleteAllProducts(context.Background(), &product.DeleteAllProductsRequest{
		UserId: userId,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": resp.Status})
}

func (a *AdminHandler) RequestGetUsersHandler(c echo.Context) error {

	resp, err := a.GrpcClientAuthService.GetUser(context.Background(), &auth.GetUserRequest{
		UserId: "ok",
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "error"})
	}

	return c.JSON(http.StatusOK, resp.GetUsers())
}

func (a *AdminHandler) RequestBanUserHandler(c echo.Context) error {
	userId := c.Param("id")
	var req BanRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "error"})
	}
	if req.Banned == true {
		respStatus, err := a.BanUser(userId, c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"status": "error"})
		}
		return c.JSON(http.StatusOK, respStatus)
	} else {
		respStatus, err := a.UnBanUser(userId, c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"status": "error"})
		}
		return c.JSON(http.StatusOK, respStatus)
	}

}

func (a *AdminHandler) RequestChangeRoleHandler(c echo.Context) error {
	userId := c.Param("id")
	var req ChangeRoleRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "error"})
	}

	resp, err := a.GrpcClientAuthService.ChangeRole(context.Background(), &auth.ChangeRoleRequest{
		UserId: userId,
		Role:   req.Role,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "error"})
	}
	if resp.Status != "ok" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "error"})
	}
	return c.JSON(http.StatusCreated, resp.Status)
}

func (a *AdminHandler) RequestBanUserByProductIdHandler(c echo.Context) error {
	productId := c.Param("id")

	url := fmt.Sprintf("http://%s/api/give-id/by-product-id/%s", config.ServiceConfig.PRODUCT_SERVICE_URL, productId)
	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "error"})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	respStatus, err := a.BanUser(string(body), c)

	return c.JSON(http.StatusCreated, respStatus)
}

func (a *AdminHandler) BanUser(userId string, c echo.Context) (string, error) {
	resp, err := a.GrpcClientAuthService.BanUser(context.Background(), &auth.BanUserRequest{
		UserId: userId,
	})
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}

func (a *AdminHandler) UnBanUser(userId string, c echo.Context) (string, error) {
	resp, err := a.GrpcClientAuthService.UnBanUser(context.Background(), &auth.UnBanUserRequest{
		UserId: userId,
	})
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}
