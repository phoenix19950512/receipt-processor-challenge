package main

import (
	"encoding/json"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var receipts = make(map[string]Receipt)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type ProcessResponse struct {
	ID string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}

type ReceiptsResponse struct {
	Receipts []Receipt `json:"data"`
}

func main() {
	http.HandleFunc("/receipts/process", processReceipt)
	http.HandleFunc("/receipts/", getPoints)
	http.ListenAndServe(":8080", nil)
}

func processReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate receipt (basic validation)
	if receipt.Retailer == "" || receipt.PurchaseDate == "" || receipt.PurchaseTime == "" || len(receipt.Items) == 0 || receipt.Total == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Generate ID
	id := uuid.New().String()
	receipts[id] = receipt

	// Return ID
	response := ProcessResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL
	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id = strings.TrimSuffix(id, "/points")

	if id == "" {
		receiptsList := make([]map[string]interface{}, 0)
		for id, receipt := range receipts {
			receiptWithID := map[string]interface{}{
				"id":           id,
				"retailer":     receipt.Retailer,
				"purchaseDate": receipt.PurchaseDate,
				"purchaseTime": receipt.PurchaseTime,
				"items":        receipt.Items,
				"total":        receipt.Total,
			}
			receiptsList = append(receiptsList, receiptWithID)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(receiptsList)
		return
	}

	receipt, exists := receipts[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Calculate points
	points := calculatePoints(receipt)

	// Return points
	response := PointsResponse{Points: points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func calculatePoints(receipt Receipt) int {
	points := 0

	// Rule 1: Alphanumeric characters in retailer name
	alphanumeric := regexp.MustCompile(`[^a-zA-Z0-9]`)
	retailerName := alphanumeric.ReplaceAllString(receipt.Retailer, "")
	points += len(retailerName)

	// Rule 2: 50 points if total is a round dollar amount
	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == math.Floor(total) {
		points += 50
	}

	// Rule 3: 25 points if total is a multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items
	points += len(receipt.Items) / 2 * 5

	// Rule 5: Item description length multiple of 3
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	// Rule 6: 6 points if purchase day is odd
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if purchase time is between 2:00 PM and 4:00 PM
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	return points
}
