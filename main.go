package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Response struct {
	Success bool   `json:"success"`
	Result  int    `json:"result"`
	Status  int    `json:"status"`
	Wallet  Wallet `json:"wallet"`
}

type Wallet struct {
	Available float64 `json:"available"`
	Hold      float64 `json:"hold"`
	Reserve   float64 `json:"reserve"`
	Currency  string  `json:"currency"`
}

// Send a GET request to the specified URL
// And return wallet.available parameter from request response
func getWalletAvailable(merchantKey string) (float64, error) {
	// Create a new request using http
	req, err := http.NewRequest("GET", "https://business.paycos.com/api/v1/balance?currency=BRL", nil)
	if err != nil {
		return 0, err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", "Bearer "+merchantKey)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	fmt.Println(resp.Status)
	if err != nil {
		return 0, err
	}

	return response.Wallet.Available, nil
}

// Run http server with "check" endpoint
func server() {
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		min, err := strconv.ParseFloat(r.URL.Query().Get("min"), 64)
		available, err := getWalletAvailable("e33ba1e2b9352bf29e16")
		if err != nil {
			panic(err)
		}

		if available < min*100 {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Insufficient funds"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Sufficient funds"))
		}
	})

	if err := http.ListenAndServe(":8001", nil); err != nil {
		panic(err)
	}
}

func main() {
	server()
}
