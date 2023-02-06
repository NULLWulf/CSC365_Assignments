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
	Categories  string  `json:"categories"`
}

type Review struct {
	ReviewID   string `json:"review_id"`
	UserID     string `json:"user_id"`
	BusinessID string `json:"business_id"`
	Stars      int    `json:"stars"`
	Useful     int    `json:"useful"`
	Funny      int    `json:"funny"`
	Cool       int    `json:"cool"`
	Text       string `json:"text"`
}

func ReadJsonFile(filePath string) (*interface{}, error) {
	var jsonData interface{}
	file, err := os.Open(filePath)
	if err != nil {
		return &jsonData, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	err = json.NewDecoder(file).Decode(&jsonData)
	if err != nil {
		return &jsonData, err
	}
	return &jsonData, nil
}
