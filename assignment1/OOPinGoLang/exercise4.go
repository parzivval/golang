package main

import (
	"encoding/json"
	"fmt"
)

type Product struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func encodeProductToJSON(p Product) string {
	jsonData, _ := json.Marshal(p)
	return string(jsonData)
}

func decodeJSONToProduct(jsonStr string) Product {
	var p Product
	_ = json.Unmarshal([]byte(jsonStr), &p)
	return p
}

func main() {
	product := Product{Name: "Laptop", Price: 999.99, Quantity: 5}

	jsonProduct := encodeProductToJSON(product)
	fmt.Println("Encoded JSON:", jsonProduct)

	decodedProduct := decodeJSONToProduct(jsonProduct)
	fmt.Printf("Decoded Product: %+v\n", decodedProduct)
}
