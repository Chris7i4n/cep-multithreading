package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func main() {
	cep := "60864260"
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		address_data, err := fetchViaCep(cep)

		if err != nil {
			log.Printf("In: %v\n", err)
			return
		}
		ch1 <- fmt.Sprintf("%s, %s, %s", address_data.Logradouro, address_data.Bairro, address_data.Localidade)
	}()
	go func() {
		address_data, err := fetchBrasilApi(cep)

		if err != nil {
			log.Printf("In: %v\n", err)
			return
		}
		ch2 <- fmt.Sprintf("%s, %s, %s", address_data.Street, address_data.Neighborhood, address_data.City)
	}()

	select {
	case response := <-ch1:
		fmt.Println("Resposta recebida da api ViaCEP:", response)
	case response := <-ch2:
		fmt.Println("Resposta recebida da api Brasilapi:", response)
	case <-time.After(1 * time.Second):
		fmt.Println("Erro: timeout de 1 segundo excedido")
	}
}

func fetchViaCep(cep string) (ViaCEPResponse, error) {
	var apiResponse ViaCEPResponse

	url := "http://viacep.com.br/ws/" + cep + "/json/"

	resp, err := http.Get(url)

	if err != nil {
		return apiResponse, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return apiResponse, fmt.Errorf("Erro na resposta: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)

	if err != nil {
		return apiResponse, err
	}

	return apiResponse, nil
}

func fetchBrasilApi(cep string) (BrasilAPIResponse, error) {
	var apiResponse BrasilAPIResponse

	url := "https://brasilapi.com.br/api/cep/v1/" + cep

	resp, err := http.Get(url)

	if err != nil {
		return apiResponse, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)

	if err != nil {
		return apiResponse, err
	}

	return apiResponse, nil
}
