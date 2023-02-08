package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Business struct {
	BusinessID  string  `json:"business_id"`
	Name        string  `json:"name"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Stars       float32 `json:"stars"`
	ReviewCount int     `json:"review_count"`
	IsOpen      int     `json:"is_open"`
	Categories  string  `json:"categories"`
}

type Review struct {
	ReviewID   string  `json:"review_id"`
	UserID     string  `json:"user_id"`
	BusinessID string  `json:"business_id"`
	Stars      float32 `json:"stars"`
	Text       string  `json:"text"`
}

// Businesses Initialize an array to store the businesses
var Businesses []Business
var Reviews []Review

func readBusinessesJson() {
	fmt.Println("Loading Business JSON data...")
	// Read the file containing business information
	file, err := os.ReadFile(businessPath)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through the file, decoding each business object and adding it to the array
	for i := 0; i < len(file); {
		var business Business
		j := i
		for j < len(file) && file[j] != '\n' {
			j++
		}
		err := json.Unmarshal(file[i:j], &business)
		if err != nil {
			log.Fatal(err)
		}
		Businesses = append(Businesses, business)
		i = j + 1
	}
	fmt.Println("Businesses Loaded: ", len(Businesses))
}

func readReviewsJson() {
	// Read the file containing business information
	fmt.Println("Loading Review JSON data...")
	file, err := os.ReadFile(reviewJsonPath)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through the file, decoding each business object and adding it to the array
	//for i := 0; i < 1000; {
	for i := 0; i < len(file); {
		var review Review
		j := i
		for j < len(file) && file[j] != '\n' {
			j++
		}
		err := json.Unmarshal(file[i:j], &review)
		if err != nil {
			fmt.Println("Review ignored: ", err)
		} else {
			Reviews = append(Reviews, review)
		}
		i = j + 1
	}
	fmt.Println("Reviews Loaded: ", len(Reviews))
}

func readReviewsJsonScannner() {
	file, err := os.Open(reviewJsonPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {

		var review Review
		err := json.Unmarshal(scanner.Bytes(), &review)
		if err != nil {
			fmt.Println("Review ignored: ", err)
		} else {
			Reviews = append(Reviews, review)
		}
		i++
		if i == 1000 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	fmt.Println("Reviews Loaded: ", len(Reviews))

}
