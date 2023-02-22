package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/bbalet/stopwords"
)

type Business struct {
	BusinessID       string             `json:"business_id"`
	Name             string             `json:"name"`
	City             string             `json:"city"`
	State            string             `json:"state"`
	Stars            float32            `json:"stars"`
	ReviewCount      int                `json:"review_count"`
	IsOpen           int                `json:"is_open"`
	Categories       string             `json:"categories"`
	CategoriesArr    []string           `json:"categories_arr" nil:"true"`
	ReviewTermsCount map[string]int     `json:"review_terms_count"`
	TermCountTotal   int                `json:"term_count_total"`
	termFrequency    map[string]float32 `json:"term_frequency"`

	TfIdf map[string]float32 `json:"tf_idf"`
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
var TermDocumentFrequency map[string]int
var ReviewTotal = 25000

func readBusinessesJson() {
	log.Println("Loading Business JSON data...")
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
			log.Println("Business ignored: ", err)
		} else {
			if business.IsOpen != 0 && strings.Contains(business.Categories, "Restaurants") && business.ReviewCount > 100 {
				business.CategoriesArr = strings.Split(business.Categories, ", ")
				business.ReviewTermsCount = make(map[string]int)
				Businesses = append(Businesses, business)
				for _, category := range business.CategoriesArr {
					categoryFrequencyTable[category]++
				}
			}
		}
		i = j + 1
	}
	log.Println("Businesses Loaded: ", len(Businesses))
	var categories []string
	for category := range categoryFrequencyTable {
		categories = append(categories, category)
	}
	sort.Slice(categories, func(i, j int) bool {
		return categoryFrequencyTable[categories[i]] > categoryFrequencyTable[categories[j]]
	})

}

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
		} else {
			for i, b := range Businesses {
				if review.BusinessID == b.BusinessID {
					tTerms := strings.Split(stopwords.CleanString(review.Text, "en", true), " ")
					for _, term := range tTerms {
						Businesses[i].ReviewTermsCount[term]++
					}
					t++
					break
				}
			}
		}
		if t == ReviewTotal {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading file:", err)
	}

	log.Printf("Reviews Loading: %d.  Businesses before removal of nulls: %d", ReviewTotal, len(Businesses))
}

func saveBusinessAsJsonArray() {
	file, err := os.Create("businesses.json")
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	err = enc.Encode(Businesses)
	if err != nil {
		log.Println("Error encoding json:", err)
	}
}

func (b Business) ToJson() []byte {
	businessJson, _ := json.Marshal(b)
	return businessJson
}
