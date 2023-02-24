package main

import (
	"bufio"
	"encoding/json"
	"github.com/bbalet/stopwords"
	"log"
	"os"
	"strings"
)

// Reads Businesses JSON data
func readBusinessesJson() {
	log.Println("Loading Business JSON data...")
	Businesses = make(map[string]Business) // Instantiate global businesses map
	// Read the file containing business information
	file, err := os.ReadFile(businessPath) // Reads entire file to file object
	if err != nil {
		log.Fatal(err)
	}
	// Loop through the file object, decoding each business object and adding it to the array
	BusinessKeyMap = make([]string, 0)
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
			if business.IsOpen != 0 && strings.Contains(business.Categories, "Restaurants") && business.ReviewCount > ReviewCount {
				business.CategoriesArr = strings.Split(business.Categories, ", ")
				business.ReviewTermsCount = make(map[string]int)
				business.ReviewTermsCountHM = *NewHashMap()
				Businesses[business.BusinessID] = business
				BusinessKeyMap = append(BusinessKeyMap, business.BusinessID)
			}
		}
		i = j + 1
	}
	log.Println("Businesses Loaded: ", len(Businesses))

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

		//for _, b := range Businesses {
		//	if b.BusinessID == review.BusinessID {
		//		// Remove stop words and split by spaces
		//		tTerms := strings.Split(stopwords.CleanString(review.Text, "en", true), " ")
		//		for _, term := range tTerms {
		//			ptr := Businesses[b.BusinessID]
		//			ptr.ReviewTermsCountHM.Add(term, 1) // Add the term to the function
		//		}
		//		t++
		//		break
		//	}
		//}

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
