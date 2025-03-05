package handlers

import (
	"ByteShop/generated/order"
	"PaymentService/internal/models"
	"PaymentService/internal/paypal"
	"PaymentService/internal/services"
	"PaymentService/pkg/database"
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type PaymentHandler struct {
	GrpcClientPaymentService order.OrderServiceClient
}

func (p *PaymentHandler) CreatePayment(c echo.Context) error {

	var req struct {
		BuyerID   string  `json:"buyerID"`
		ProductID string  `json:"productID"`
		Amount    float64 `json:"amount"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверные данные"})
	}

	amount, err := services.GetAmount(req.ProductID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Не удалось получить цену товара"})
	}

	req.Amount = amount
	sellerID, err := services.GetSellerId(req.ProductID)

	if err != nil || sellerID == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка поиска ID"})
	}

	sellerPayPalEmail, err := services.GetSellerPayPalEmail(sellerID)
	log.Println(sellerPayPalEmail)
	if err != nil || sellerPayPalEmail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Продавец не указал PayPal-аккаунт"})
	}

	paypalOrderID, approveURL, err := paypal.CreatePayPalOrder(req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка создания PayPal заказа"})
	}

	payment := models.Payment{
		ProductID:     req.ProductID,
		BuyerID:       req.BuyerID,
		SellerID:      sellerID,
		Amount:        req.Amount,
		Fee:           req.Amount * 0.10,
		NetAmount:     req.Amount * 0.90,
		Status:        "pending",
		PaymentMethod: "paypal",
		TransactionID: paypalOrderID,
	}

	if err := database.DB.Create(&payment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка записи в БД"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"payment_id":   payment.ID,
		"paypal_order": paypalOrderID,
		"approve_url":  approveURL,
	})
}

func (p *PaymentHandler) CapturePayment(c echo.Context) error {
	var req struct {
		PayPalOrder string `json:"paypal_order_id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверные данные"})
	}

	var payment models.Payment

	if err := database.DB.Where("transaction_id = ?", req.PayPalOrder).First(&payment).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Платеж не найден"})
	}

	if payment.Status != "pending" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Платеж уже обработан"})
	}

	transactionID, err := paypal.CapturePayPalOrder(req.PayPalOrder)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка подтверждения PayPal платежа"})
	}
	log.Println("Запрос к GetOrder:", payment.BuyerID, payment.ProductID, payment.Amount)

	resp, err := p.GrpcClientPaymentService.GetOrder(context.TODO(), &order.GetOrderRequest{
		BuyerID:   payment.BuyerID,
		ProductID: payment.ProductID,
		Amount:    payment.Amount,
	})
	log.Println(resp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "error"})
	}

	payment.Status = "completed"
	payment.TransactionID = transactionID

	if err := database.DB.Save(&payment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка обновления БД"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Оплата подтверждена"})
}
