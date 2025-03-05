package services

import (
	"PaymentService/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetSellerPayPalEmail(sellerID string) (string, error) {

	url := fmt.Sprintf("http://%s/api/user/%s/paypal", config.ServiceConfig.AUTH_SERVICE_URL, sellerID)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("не удалось получить PayPal email")
	}

	var response struct {
		PayPalEmail string `json:"payPalEmail"`
		Status      string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.PayPalEmail, nil
}

func GetSellerId(productID string) (string, error) {

	url := fmt.Sprintf("http://%s/api/give-id/by-product-id/%s", config.ServiceConfig.PRODUCT_SERVICE_URL, productID)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("не удалось получить sellerID")
	}

	var response struct {
		UserID string `json:"userID"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	return response.UserID, nil
}

func GetAmount(productID string) (float64, error) {
	url := fmt.Sprintf("http://%s/api/give-price/by-product-id/%s", config.ServiceConfig.PRODUCT_SERVICE_URL, productID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("сервер вернул ошибку: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	Response := struct {
		Amount float64 `json:"amount"`
	}{}

	if err := json.Unmarshal(body, &Response); err != nil {
		return 0, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	log.Println("Цена товара:", Response.Amount)
	return Response.Amount, nil
}
