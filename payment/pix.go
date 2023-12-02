package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// PaymentResponse estrutura para armazenar a resposta da API de pagamento
type PaymentResponse struct {
	QRCode string `json:"qr_code"`
	// Outros campos da resposta que você possa precisar
}

// PIX realiza um pagamento PIX e retorna a resposta da API
func PIX(valor int, mpToken string, ID int) (*PaymentResponse, error) {
	url := "https://api.mercadopago.com/v1/payments"

	payload := map[string]interface{}{
		"transaction_amount": valor,
		"payment_method_id":  "pix",
		"payer": map[string]interface{}{
			"email": "inijbe@gmail.com",
		},
	}

	

	jsonValue, _ := json.Marshal(payload)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("Erro ao criar a requisição:", err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Idempotency-Key", fmt.Sprintf("%v", ID))
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

	// Verifica se a requisição foi bem-sucedida (código de status 200)
	if response.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Erro ao analisar o JSON da resposta:", err)
			return nil, err
		}

		// Verifica se a chave "point_of_interaction" está presente no JSON
		if poi, ok := result["point_of_interaction"].(map[string]interface{}); ok {
			// Verifica se a chave "qr_code" está presente dentro de "point_of_interaction"
			if qrCode, ok := poi["transaction_data"].(map[string]interface{})["qr_code"]; ok {
				fmt.Println("QR Code:", qrCode)
				paymentResponse := &PaymentResponse{
					QRCode: qrCode.(string),
					// Defina outros campos necessários do pagamento aqui
				}
				return paymentResponse, nil
			}
		}
		fmt.Println("Campo 'qr_code' não encontrado no JSON.")
	} else {
		fmt.Println("Erro na requisição:", response.Status)
		fmt.Println("Corpo da resposta:", string(body))
		return nil, errors.New(fmt.Sprintf("Erro na requisição: %s", response.Status))
	}
	return nil, nil
}

// GetQRCode realiza o pagamento PIX e retorna o QR Code
func GetQRCode(valor int, mpToken string) (string, error) {
	paymentResponse, err := PIX(valor, mpToken, 0)
	if err != nil {
		return "", err
	}

	if paymentResponse != nil {
		return paymentResponse.QRCode, nil
	}

	return "", errors.New("Resposta de pagamento vazia")
}
