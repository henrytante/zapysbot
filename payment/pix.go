package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"telebotgo/payment/status"
	"math/rand"
	"time"
)

type PaymentResponse struct {
	QRCode string `json:"qr_code"`
	ID     int    `json:"id"`
}

func generateRandomIdempotency(length int) string {
	rand.Seed(time.Now().UnixNano())

	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)

	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charSet))
		code[i] = charSet[randomIndex]
	}

	return string(code)
}

func PIX(valor float64, mpToken string, ID int, ch chan<- bool) (*PaymentResponse, error) {
	url := "https://api.mercadopago.com/v1/payments"

	payload := map[string]interface{}{
		"transaction_amount": valor,
		"payment_method_id":  "pix",
		"payer": map[string]interface{}{
			"email": "inijbe@gmail.com",
		},
	}
	expirationTime := time.Now().Add(8 * time.Minute)
	expirationFormatted := expirationTime.Format("2006-01-02T15:04:05.000-07:00")
	payload["date_of_expiration"] = expirationFormatted

	jsonValue, _ := json.Marshal(payload)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("Erro ao criar a requisição:", err)
		return nil, err
	}
	idempotency := generateRandomIdempotency(8)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Idempotency-Key", idempotency)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mpToken))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Erro ao enviar a requisição:", err)
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erro ao ler corpo da resposta:", err)
		return nil, err
	}

	if response.StatusCode == http.StatusCreated {
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Erro ao analisar o JSON da resposta:", err)
			return nil, err
		}

		if id, ok := result["id"].(float64); ok {
			fmt.Println("PIX gerado, ID :", int(id))
			go status.Status(mpToken, int(id), ch)
			var qrCode string

			if poi, ok := result["point_of_interaction"].(map[string]interface{}); ok {
				if qr, ok := poi["transaction_data"].(map[string]interface{})["qr_code"]; ok {
					fmt.Println("QR Code:", qr)
					qrCode = qr.(string)
				}
			}

			paymentResponse := &PaymentResponse{
				QRCode: qrCode,
				ID:     int(id),
			}
			return paymentResponse, nil
		}
		fmt.Println("Campo 'id' não encontrado no JSON.")
	} else {
		fmt.Println("Erro na requisição:", response.Status)
		fmt.Println("Corpo da resposta:", string(body))
		return nil, errors.New(fmt.Sprintf("Erro na requisição: %s", response.Status))
	}
	return nil, nil
}

func GetQRCode(valor float64, mpToken string) (string, error) {
	ch := make(chan bool)
	paymentResponse, err := PIX(valor, mpToken, 0, ch) // Adicionado o canal como argumento
	if err != nil {
		return "", err
	}

	if paymentResponse != nil {
		return paymentResponse.QRCode, nil
	}

	return "", errors.New("Resposta de pagamento vazia")
}


func GetPixInfo(mpToken string, pixID int) (*PaymentResponse, error) {
	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%d", pixID)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Erro ao criar a requisição:", err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mpToken))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Erro ao enviar a requisição:", err)
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erro ao ler corpo da resposta:", err)
		return nil, err
	}

	if response.StatusCode == http.StatusCreated {
		var result PaymentResponse
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Erro ao analisar o JSON da resposta:", err)
			return nil, err
		}

		return &result, nil
	}

	fmt.Println("Erro na requisição:", response.Status)
	fmt.Println("Corpo da resposta:", string(body))
	return nil, errors.New(fmt.Sprintf("Erro na requisição: %s", response.Status))
}
