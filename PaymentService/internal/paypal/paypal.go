package paypal

import (
	"PaymentService/config"
	"PaymentService/internal/services"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type PayPalClient struct {
	ClientID     string
	ClientSecret string
	APIBaseURL   string
}

// var loadConfig = config.LoadConfig()
var paypal = &PayPalClient{
	ClientID:     config.ServiceConfig.PAYPAL_CLIENT_ID,
	ClientSecret: config.ServiceConfig.PAYPAL_SECRET,
	APIBaseURL:   config.ServiceConfig.PAYPAL_API,
}

func CreatePayPalOrder(amount float64) (string, string, error) {

	requestBody, _ := json.Marshal(map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"amount": map[string]interface{}{
					"currency_code": "USD",
					"value":         amount,
				},
			},
		},
	})

	req, _ := http.NewRequest("POST", paypal.APIBaseURL+"/v2/checkout/orders", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(paypal.ClientID, paypal.ClientSecret)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)

	if res.StatusCode != http.StatusCreated {
		return "", "", errors.New("не удалось создать заказ")
	}

	orderID := response["id"].(string)
	approveURL := ""
	for _, link := range response["links"].([]interface{}) {
		linkMap := link.(map[string]interface{})
		if linkMap["rel"] == "approve" {
			approveURL = linkMap["href"].(string)
			break
		}
	}

	return orderID, approveURL, nil
}

func CapturePayPalOrder(orderID string) (string, error) {

	reqURL := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", paypal.APIBaseURL, orderID)

	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(paypal.ClientID, paypal.ClientSecret)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return "", errors.New("не удалось подтвердить оплату")
	}

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)

	transactionID := ""
	if purchaseUnits, ok := response["purchase_units"].([]interface{}); ok {
		if len(purchaseUnits) > 0 {
			payments := purchaseUnits[0].(map[string]interface{})["payments"].(map[string]interface{})
			captures := payments["captures"].([]interface{})
			if len(captures) > 0 {
				transactionID = captures[0].(map[string]interface{})["id"].(string)
			}
		}
	}

	if transactionID == "" {
		return "", errors.New("не удалось получить transaction ID")
	}

	return transactionID, nil
}

type PayoutRequest struct {
	SenderBatchHeader struct {
		SenderBatchID string `json:"sender_batch_id"`
		EmailSubject  string `json:"email_subject"`
	} `json:"sender_batch_header"`
	Items []PayoutItem `json:"items"`
}

type PayoutItem struct {
	RecipientType string `json:"recipient_type"`
	Receiver      string `json:"receiver"`
	Amount        struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Note         string `json:"note"`
	SenderItemID string `json:"sender_item_id"`
}

func PayoutToSeller(amount float64, sellerID string) error {

	sellerEmail, err := services.GetSellerPayPalEmail(sellerID)
	if err != nil {
		return fmt.Errorf("ошибка получения PayPal email: %w", err)
	}

	payoutRequest := PayoutRequest{}
	payoutRequest.SenderBatchHeader.SenderBatchID = fmt.Sprintf("payout_%d", sellerID)
	payoutRequest.SenderBatchHeader.EmailSubject = "Вам отправлен платеж!"
	payoutRequest.Items = []PayoutItem{
		{
			RecipientType: "EMAIL",
			Receiver:      sellerEmail,
			Amount: struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			}{
				Value:    fmt.Sprintf("%.2f", amount),
				Currency: "USD",
			},
			Note:         "Выплата за проданный товар",
			SenderItemID: fmt.Sprintf("payout_item_%d", sellerID),
		},
	}

	payload, err := json.Marshal(payoutRequest)
	if err != nil {
		return err
	}

	reqURL := fmt.Sprintf("%s/v1/payments/payouts", paypal.APIBaseURL)
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(paypal.ClientID, paypal.ClientSecret)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return errors.New("ошибка отправки выплаты")
	}

	return nil
}
