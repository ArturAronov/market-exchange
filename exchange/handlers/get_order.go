package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"market-exchange/utils"
)

func GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to recieve order: Expexted method GET, recieved" + r.Method))
	}

	order, orderErr := utils.OrderUriParser(r.RequestURI)
	if orderErr != nil {
		w.WriteHeader(http.StatusBadRequest)

		orderMapErr := map[string]string{"error": orderErr.Error()}
		responseJsonErr := json.NewEncoder(w).Encode(orderMapErr)
		if responseJsonErr != nil {
			log.Printf("Failed to encode JSON: %v", responseJsonErr)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	fmt.Println(order)
}
