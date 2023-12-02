package status

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type PaymentStatus struct {
	Status string `json:"status"`
	// Adicione outros campos conforme necessário
}

func checkPaymentStatus(paymentID int, mpToken string, ch chan<- bool) {
	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%v", paymentID)
	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Erro ao criar requisição:", err)
			ch <- false
			return
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mpToken))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Erro ao enviar requisição:", err)
			ch <- false
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Erro ao ler resposta:", err)
			ch <- false
			return
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Erro na resposta da API:", url, resp.Status)
			ch <- false
			return
		}

		var paymentStatus PaymentStatus
		if err := json.Unmarshal(body, &paymentStatus); err != nil {
			fmt.Println("Erro ao decodificar JSON:", err)
			ch <- false
			return
		}

		fmt.Printf("Pagamento de %v: %v\n", paymentID, paymentStatus.Status)

		if paymentStatus.Status == "approved" {
			fmt.Println("Pagamento aprovado!")
			ch <- true
			return
		}

		time.Sleep(5 * time.Second)
	}
}

func Status(mpToken string, ID int, ch chan<- bool) {
	go checkPaymentStatus(ID, mpToken, ch)
}
