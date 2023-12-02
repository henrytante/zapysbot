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

func checkPaymentStatus(paymentID int, mpToken string) {
	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%v", paymentID)
	for {
		// Criar uma requisição GET para a API do Mercado Pago para verificar o status do pagamento
		req, err := http.NewRequest("GET", fmt.Sprintf("%s", url), nil)
		if err != nil {
			fmt.Println("Erro ao criar requisição:", err)
			return
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mpToken))

		// Enviar a requisição HTTP
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Erro ao enviar requisição:", err)
			return
		}

		// Ler o corpo da resposta
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Erro ao ler resposta:", err)
			return
		}

		// Verificar o status da resposta
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Erro na resposta da API:", url,resp.Status)
			return
		}

		// Decodificar o JSON retornado pela API para a estrutura de dados desejada
		var paymentStatus PaymentStatus
		if err := json.Unmarshal(body, &paymentStatus); err != nil {
			fmt.Println("Erro ao decodificar JSON:", err)
			return
		}

		// Acessar e exibir o campo "status" da resposta JSON
		
		fmt.Printf("Pagamento de %v: ", paymentID, paymentStatus.Status)

		// Verificar se o pagamento foi aprovado
		if paymentStatus.Status == "approved" {
			fmt.Println("Pagamento aprovado!")
			// Adicione aqui a lógica para o que deseja fazer quando o pagamento for aprovado
			break
		}

		// Esperar alguns segundos antes de verificar novamente
		time.Sleep(5 * time.Second)
	}
}

func Status(mpToken string, ID int) {
	
	// Chame a função para verificar o status do pagamento de forma assíncrona
	go checkPaymentStatus(ID, mpToken)

	// Mantenha o programa rodando para continuar verificando o status do pagamento em segundo plano
	select {}
}
