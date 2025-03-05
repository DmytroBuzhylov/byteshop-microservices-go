package handlers

import (
	"ProductService/config"
	"ProductService/internal/models"
	"ProductService/internal/services"
	"ProductService/pkg/database"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type UserHandler struct {
}

func (h *UserHandler) CreateProduct(c echo.Context) error {
	userID := c.FormValue("user_id")

	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	name := c.FormValue("product-name")
	description := c.FormValue("product-features")
	price, _ := strconv.ParseFloat(c.FormValue("product-price"), 64)
	deliveryType := c.FormValue("delivery-type")
	//isUnlimited := c.FormValue("is_unlimited") == "true"

	var link string
	var keys []string

	if deliveryType == "key" {
		items := c.FormValue("product-keys")
		keys = strings.Split(items, "\n")
	} else if deliveryType == "link" {
		link = c.FormValue("product-link")
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Такого типа не существует"})
	}

	category := c.FormValue("category")
	if !services.CheckCategory(category) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Такой категории не существует"})
	}

	file, err := c.FormFile("product-image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Фото товара обязательно"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при обработке файла"})
	}
	defer src.Close()

	filePath := "uploads/" + file.Filename
	dst, err := os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при сохранении файла"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при копировании файла"})
	}

	product := models.Product{
		ID:          uuid.New(),
		UserID:      uuidUserID,
		Name:        name,
		Description: description,
		Price:       price,
		ImageURL:    "/uploads/" + file.Filename,
		Category:    category,
		Type:        deliveryType,
	}

	if deliveryType == "key" {
		for _, value := range keys {
			productItem := models.ProductItems{
				ID:        uuid.New(),
				ProductID: product.ID,
				Value:     value,
				IsSolid:   false,
				SoldAt:    nil,
			}
			err = database.DB.Create(&productItem).Error
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при сохранении товара"})
			}
		}
	} else if deliveryType == "link" {
		product.Link = link
	}

	if err := database.DB.Create(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при сохранении товара"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Товар успешно добавлен"})
}

func (h *UserHandler) UpdateProduct(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]string{"status": "ok"})
}

func (h *UserHandler) GetProducts(c echo.Context) error {
	var products []models.Product
	err := database.DB.Model(&models.Product{}).Select("*").Find(&products).Order("created_at DESC").Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "error"})
	}

	return c.JSON(http.StatusCreated, products)
}

func (h *UserHandler) GetSellerProducts(c echo.Context) error {
	header := c.Request().Header.Get("Authorization")
	if header == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"status": "missing token"})
	}

	url := fmt.Sprintf("http://%s/api/give-user-id", config.ServiceConfig.AUTH_SERVICE_URL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "request error"})
	}

	req.Header.Set("Authorization", header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "request failed"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(resp.StatusCode, map[string]string{"status": "invalid response"})
	}

	var result struct {
		UserID string `json:"userID"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "json decode error"})
	}

	var products []models.Product

	err = database.DB.Table("products").Where("user_id = ?", result.UserID).Find(&products).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "db error"})
	}

	return c.JSON(http.StatusOK, &products)
}

func (h *UserHandler) GetProductsCategory(c echo.Context) error {
	category := c.Param("category")
	products, err := h.GetProductsByCategory(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка получения товаров"})
	}

	return c.JSON(http.StatusOK, products)
}

func (h *UserHandler) GetProductsByCategory(category string) ([]models.Product, error) {
	var productsCatagory []models.Product
	err := database.DB.Model(&models.Product{}).Where("category = ?", category).Find(&productsCatagory).Error
	if err != nil {
		return nil, err
	}
	return productsCatagory, nil
}

func (h *UserHandler) GetProductById(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	err := database.DB.Model(&models.Product{}).Where("id = ?", id).Find(&product).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка получения товара"})
	}
	return c.JSON(http.StatusOK, product)
}

func (h *UserHandler) GetUserByProductId(c echo.Context) error {

	id := c.Param("id")
	var product models.Product
	err := database.DB.Model(&models.Product{}).Where("id = ?", id).Find(&product).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка получения товара"})
	}
	return c.JSON(http.StatusOK, map[string]uuid.UUID{"userID": product.UserID})
}

func (h *UserHandler) GetPriceByProductId(c echo.Context) error {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректный UUID"})
	}

	var product models.Product
	err = database.DB.Where("id = ?", productID).First(&product).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка получения товара"})
	}

	return c.JSON(http.StatusOK, map[string]float64{"amount": product.Price})
}

func (h *UserHandler) GetAllProducts(c echo.Context) error {
	id := c.Param("userID")
	var products []models.Product
	result := database.DB.Where("user_id = ?", id).Find(&products)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка получения товаров"})
	}
	return c.JSON(http.StatusOK, products)
}

func (h *UserHandler) GetAllProductId(c echo.Context) error {
	id := c.Param("userID")
	var productIDs []string
	result := database.DB.Model(&models.Product{}).Where("user_id = ?", id).Pluck("id", &productIDs)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка получения товаров"})
	}
	return c.JSON(http.StatusOK, productIDs)
}
