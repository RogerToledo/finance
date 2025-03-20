package controller

import (
	"encoding/json"
	"net/http"
)

func HTTPResponse(w http.ResponseWriter, message any, statusCode int) map[string]any {
	resp := map[string]any{
		"StatusCode": statusCode,
		"Message": message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // ou "*" para todas as origens
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Métodos permitidos
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	json.NewEncoder(w).Encode(resp)

	return resp
}
