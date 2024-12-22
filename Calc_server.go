package main

import (
"encoding/json"
"fmt"
"net/http"
)

type CalculationRequest struct {
	Expression string json:"expression"
}

type CalculationResponse struct {
	Result string json:"result,omitempty"
	Error  string json:"error,omitempty"
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	result, err := Calc(req.Expression)
	response := CalculationResponse{}

	if err != nil {
		if err.Error() == "число не было найдено" || err.Error() == "некорректное выражение" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			response.Error = "Expression is not valid"
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			response.Error = "Internal server error"
		}
	} else {
		response.Result = fmt.Sprintf("%f", result)
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
