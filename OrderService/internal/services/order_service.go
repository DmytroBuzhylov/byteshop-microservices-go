package services

import (
	"OrderService/config"
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type EmailData struct {
	ProductName  string
	ProductKey   string
	OrderNumber  string
	PurchaseDate string
}

type EmailResponse struct {
	Email string `json:"email"`
}

func SendMail(to string, data EmailData, flag string) error {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal("Ошибка получения рабочего каталога:", err)
	}

	var templatePath string
	if flag == "key" {
		templatePath = filepath.Join(basePath, "..", "OrderService", "sample", "email_template_key.html")
	} else if flag == "link" {
		templatePath = filepath.Join(basePath, "..", "OrderService", "sample", "email_template_link.html")
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal("Ошибка загрузки шаблона:", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("ошибка генерации HTML: %w", err)
	}
	m := gomail.NewMessage()
	m.SetHeader("From", config.ServiceConfig.EMAIL)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Ваш заказ успешно оформлен")
	m.SetBody("text/html", body.String())
	d := gomail.NewDialer("localhost", 1025, "", "")
	if err = d.DialAndSend(m); err != nil {
		return fmt.Errorf("ошибка отправки письма: %w", err)
	}

	fmt.Println("Письмо отправлено!")
	return nil
}

func GetUserEmail(userID string) (string, error) {
	url := fmt.Sprintf("http://%s/api/email/%s", config.ServiceConfig.AUTH_SERVICE_URL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка запроса: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("сервер вернул ошибку: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения тела ответа: %w", err)
	}

	var result EmailResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	return result.Email, nil
}
