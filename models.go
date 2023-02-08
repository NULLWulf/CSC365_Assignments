package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Business struct {
	BusinessID    string   `json:"business_id"`
	Name          string   `json:"name"`
	City          string   `json:"city"`
	State         string   `json:"state"`
	Stars         float32  `json:"stars"`
	ReviewCount   int      `json:"review_count"`
	IsOpen        int      `json:"is_open"`
	Categories    string   `json:"categories"`
	CategoriesArr []string `json:"categories_arr" nil:"true"`
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
var BusinessIdList []string

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
			fmt.Println("Business ignored: ", err)
		} else {
			if business.IsOpen != 0 && strings.Contains(business.Categories, "Restaurants") {
				business.CategoriesArr = strings.Split(business.Categories, ", ")
				Businesses = append(Businesses, business)
				BusinessIdList = append(BusinessIdList, business.BusinessID)
			}
		}
		i = j + 1
	}
	fmt.Println("Businesses Loaded: ", len(Businesses))
	fmt.Println(Businesses[1:10])
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

func delimitCategories(categories string) []string {
	return strings.Split(categories, ",")
}
