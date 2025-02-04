# Receipt Points Calculator API

[![License](https://img.shields.io/badge/license-MIT-blue.svg)]()

A robust web service that processes receipts and calculates reward points based on specific rules. This service provides a RESTful API for receipt processing and points calculation.

## Features

- Process receipts and generate unique IDs
- Calculate points based on receipt data
- Retrieve points for specific receipts
- View all processed receipts

## Prerequisites
- Go 1.16+
- Git
- curl or Postman (for API testing)

## API Endpoints

### 1. Process Receipt
- **POST** `/receipts/process`
- Submits a receipt for processing
- Returns a unique ID

### 2. Get Points
- **GET** `/receipts/{id}/points`
- Returns points calculated for a specific receipt

### 3. List All Receipts
- **GET** `/receipts`
- Returns all processed receipts

## Points Calculation Rules

Points are awarded based on these criteria:
- Retailer name: 1 point per alphanumeric character
- Round dollar amounts: 50 points
- Quarter multiples: 25 points
- Item pairs: 5 points per pair
- Item description length: Special calculation for multiples of 3
- Purchase date: 6 points for odd days
- Purchase time: 10 points between 2:00 PM and 4:00 PM

## Running the Application

1. Install Go (1.16 or later)
2. Clone the repository
3. Install dependencies:
   ```bash
   go mod init receipt-processor
   go mod tidy
   ```
4. Run the server:
   ```bash
   go run main.go
   ```
   The server will start on `http://localhost:8080`.

## Example Request
Request:  
**POST** `/receipts/process`
```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    }
  ],
  "total": "6.49"
}
```
Response:
```json
{
    "id": "7fb1377b-b223-49d9-a31a-5a02701dd310"
}
```

## Tech Stack
- Go
- Standard library HTTP server
- In-memory storage
- UUID generation for receipt IDs

## Data Storage
All data is stored in memory and will be cleared when the server restarts.

## Troubleshooting

Common issues and solutions:
- Server won't start: Check port 8080 availability
- API returns 400: Verify JSON payload format
