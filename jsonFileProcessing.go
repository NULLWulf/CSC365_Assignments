package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bbalet/stopwords"
)

// Reads Businesses JSON data
func readBusinessesJson() {
	log.Println("Loading Business JSON data...")
	businessMap := NewHashMap() // Create new hash map for businesses
	t := 0
	// Create directory for fileblock if it does not exist
	err := os.MkdirAll("fileblock", 0777)
	if err != nil {
		log.Fatal(err)
	}

	// Read the file containing business information
	file, err := os.ReadFile(businessPath) // Reads entire file to file object
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(file); {
		var business Business // temporary Business variable
		j := i
		for j < len(file) && file[j] != '\n' {
			j++
		}
		err := json.Unmarshal(file[i:j], &business)
		if err != nil {
			log.Println("Business ignored: ", err)
		} else {
			// Get businesses
			// go func() {
			if business.IsOpen != 0 && strings.Contains(business.Categories, "Restaurants") && business.ReviewCount > ReviewCount {

				replaceSlash := strings.ReplaceAll(business.Categories, "\u0026", ",")
				replaceSlash = strings.ReplaceAll(replaceSlash, "/", ", ")

				business.CategoriesArr = strings.Split(replaceSlash, ", ")

				// Write business to JSON file
				file := fmt.Sprintf("fileblock/%s.json", business.BusinessID)
				businessFile, err := os.Create(file)
				if err != nil {
					log.Fatal(err)
				}
				defer businessFile.Close()

				businessJson, err := json.Marshal(business)
				if err != nil {
					log.Fatal(err)
				}

				_, err = businessFile.Write(businessJson)
				if err != nil {
					log.Fatal(err)
				}
				businessMap.Add(business.BusinessID, file)
				t++

			}
			// }()

		}
		i = j + 1
		if t == 10000 {
			break
		}
	}
	log.Println("Businesses Loaded: ", businessMap.size)

}

// Reads through the reviews JSON file.  When a review is found for a bussiness it removes
// stop words from the text, splits the text into an array, and adds the term to a raw
// count frequency map for the review's respective business.
func readReviewsJsonScannner() {
	log.Println("Reading review JSON data...")

	file, err := os.Open(reviewJsonPath)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	t := 0
	for scanner.Scan() {
		var review Review

		err := json.Unmarshal(scanner.Bytes(), &review)
		if err != nil {
			log.Println("Review ignored: ", err)
			continue
		}

		for _, b := range Businesses {
			if b.BusinessID == review.BusinessID {
				// Remove stop words and split by spaces
				tTerms := strings.Split(stopwords.CleanString(review.Text, "en", true), " ")
				for _, term := range tTerms {
					ptr := Businesses[b.BusinessID]
					if ptr.ReviewTermsCount == nil {
						ptr.ReviewTermsCount = make(map[string]int)
					}
					ptr.ReviewTermsCount[term]++
				}
				t++
				break
			}
		}

		// Total on total possible rules that can be added to a business.
		// Can be modified globally
		if t == ReviewTotal {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading file:", err)
	}

	log.Printf("Reviews Loading: %d.  Businesses before removal of nulls: %d", ReviewTotal, len(Businesses))
}

func ReadBusinessJSON2() ([]Business, []BusinessDataPoint) {
	log.Println("Loading Business JSON data...")
	var Businesses []Business
	var BusinessDPS []BusinessDataPoint
	t := 0
	// Create directory for fileblock if it does not exist
	err := os.MkdirAll("fileblock", 0777)
	if err != nil {
		log.Fatal(err)
	}

	// Read the file containing business information
	file, err := os.ReadFile("yelp_academic_dataset_business.json") // Reads entire file to file object
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(file); {
		var business Business // temporary Business variable
		j := i
		for j < len(file) && file[j] != '\n' {
			j++
		}
		err := json.Unmarshal(file[i:j], &business)
		if err != nil {
			log.Println("Business ignored: ", err)
		} else {
			if business.IsOpen != 0 && strings.Contains(business.Categories, "Restaurants") && business.ReviewCount > 100 {
				// // Write business to JSON file
				// file := fmt.Sprintf("fileblock/%s.json", business.BusinessID)
				// businessFile, err := os.Create(file)
				// if err != nil {
				// 	log.Fatal(err)
				// }
				// defer businessFile.Close()

				// businessJson, err := json.Marshal(business)
				// if err != nil {
				// 	log.Fatal(err)
				// }

				// _, err = businessFile.Write(businessJson)
				// if err != nil {
				// 	log.Fatal(err)
				// }
				Businesses = append(Businesses, business)
				BusinessDPS = append(BusinessDPS, BusinessDataPoint{BusinessID: business.BusinessID, Latitude: business.Latitude, Longitude: business.Longtitude, ReviewScore: float32(business.Stars)})
				t++

			}
			// }()

		}
		i = j + 1
		if t == 10000 {
			break
		}
	}
	log.Println("Businesses Loaded: ", len(Businesses))

	return Businesses, BusinessDPS

}
