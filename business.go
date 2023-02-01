package main

import (
	"encoding/json"
	"os"
)

type Business struct {
	BusinessID  string  `json:"business_id"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	PostalCode  string  `json:"postal_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Stars       float64 `json:"stars"`
	ReviewCount int     `json:"review_count"`
	IsOpen      int     `json:"is_open"`
	Attributes  struct {
		ByAppointmentOnly string `json:"ByAppointmentOnly"`
	} `json:"attributes"`
	Categories string      `json:"categories"`
	Hours      interface{} `json:"hours"`
}

func ReadJSONFile(filePath string) (Business, error) {
	var business Business
	file, err := os.Open(filePath)
	if err != nil {
		return business, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&business)
	if err != nil {
		return business, err
	}
	return business, nil
}
