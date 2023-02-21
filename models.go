package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Business struct {
	BusinessID       string         `json:"business_id"`
	Name             string         `json:"name"`
	City             string         `json:"city"`
	State            string         `json:"state"`
	Stars            float32        `json:"stars"`
	ReviewCount      int            `json:"review_count"`
	IsOpen           int            `json:"is_open"`
	Categories       string         `json:"categories"`
	CategoriesArr    []string       `json:"categories_arr" nil:"true"`
	Reviews          []Review       `json:"reviews"`
	ReviewTermsCount map[string]int `json:"review_terms_count"`
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

func readBusinessesJson() {
	fmt.Println("Loading Business JSON data...")
	// Read the file containing business information
	file, err := os.ReadFile(businessPath)
	if err != nil {
		log.Fatal(err)
	}
	categoryFrequencyTable := make(map[string]int)
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
			if business.IsOpen != 0 && strings.Contains(business.Categories, "Restaurants") && business.ReviewCount > 100 {
				business.CategoriesArr = strings.Split(business.Categories, ", ")
				Businesses = append(Businesses, business)
				for _, category := range business.CategoriesArr {
					categoryFrequencyTable[category]++
				}
			}
		}
		i = j + 1
	}
	fmt.Println("Businesses Loaded: ", len(Businesses))
	var categories []string
	for category := range categoryFrequencyTable {
		categories = append(categories, category)
	}
	sort.Slice(categories, func(i, j int) bool {
		return categoryFrequencyTable[categories[i]] > categoryFrequencyTable[categories[j]]
	})

	//fmt.Println("Sorted Category Frequency:")
	for _, category := range categories {
		fmt.Println(category, ":", categoryFrequencyTable[category])
	}
}

func readReviewsJsonScannner() {
	file, err := os.Open(reviewJsonPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	t := 0
	for scanner.Scan() {

		var review Review

		err := json.Unmarshal(scanner.Bytes(), &review)
		if err != nil {
			fmt.Println("Review ignored: ", err)
		} else {
			for i, b := range Businesses {
				if review.BusinessID == b.BusinessID {
					//println("Review found for business: ", b.Name)
					b.Reviews = append(b.Reviews, review)
					Businesses[i].Reviews = append(Businesses[i].Reviews, review)
					t++
					break
				}
			}
		}
		if t == 25000 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	fmt.Println("Reviews Loading.  Businesses before: ", len(Businesses))
	RemoveNullReviewsFromBusinesses()
	fmt.Println("Reviews Loaded.  Businesses with reviews: ", len(Businesses))
}

func saveBusinessAsJsonArray() {
	file, err := os.Create("businesses.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	err = enc.Encode(Businesses)
	if err != nil {
		fmt.Println("Error encoding json:", err)
	}
}

func (b Business) ToJson() []byte {
	businessJson, _ := json.Marshal(b)
	return businessJson
}

// RemoveNullReviewsFromBusinesses removes businesses with no reviews from the Businesses array by
// creating a new array and copying over the businesses with reviews
func RemoveNullReviewsFromBusinesses() {
	var newBusinesses []Business
	for _, b := range Businesses {
		if len(b.Reviews) != 0 {
			newBusinesses = append(newBusinesses, b)
		}
	}
	Businesses = newBusinesses
}
