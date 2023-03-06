package main

import (
	"log"
	"math"
	"math/rand"
	"sort"
)

this will tasdasdhrow an error 

// RemoveNullReviewsFromBusinesses removes businesses with no reviews from the Businesses array by
// creating a new array and copying over the businesses with reviews
// It also calculates the term frequency for each term in the reviews of the business
// and increments the document frequency for each term
func removeNullReviewsCalculateFrequency() {
	newBusinesses := make(map[string]Business)
	tdf := make(map[string]int) // temporary global term frequency list for idf calculation

	for _, b := range Businesses {
		if len(b.ReviewTermsCount) != 0 {
			delete(b.ReviewTermsCount, "") // delete empty string key if it exists
			tCount := 0                    // total numbers of terms in all of the reviews for a given business
			for _, count := range b.ReviewTermsCount {
				tCount += count
			}
			b.TermFrequency = make(map[string]float32)
			for k, v := range b.ReviewTermsCount {
				// Calculate Term Frequency
				b.TermFrequency[k] = float32(v) / float32(tCount)
				// Increment Document Frequency should only be done once per document (business)
				tdf[k]++
			}
			b.TermCountTotal = tCount
			newBusinesses[b.BusinessID] = b
		}
	}

	// Assign the updated map to the global variable
	Businesses = newBusinesses
	TermDocumentFrequency = tdf
	log.Println("Businesses loaded after null removal: ", len(Businesses))
}

// Calculate the tf-idf of a term within a business relative to all of the businesses
func calculatetfIdf() {
	for i, b := range Businesses {
		b.TfIdf = make(map[string]float32)
		for k, v := range b.TermFrequency {
			// Calculate the tf-idf for a given a term for the business
			b.TfIdf[k] = v * float32(math.Log(float64(len(Businesses))/float64(TermDocumentFrequency[k])))
		}
		Businesses[i] = b

	}
	log.Printf("Businesses: %d, TermDocumentFrequencyNumbers: %d\n", len(Businesses), len(TermDocumentFrequency))
}

// Sorts tf-idf maps for businesses in decscendign order
func sortTfIdf() {
	log.Println("Sorting TF-IDF...")

	tempBizMap := make(map[string]Business)
	// Sort the TfIdf maps in each Business in descending order of value
	for _, b := range Businesses {
		tfidfSlice := make([]mapEntry, 0, len(b.TfIdf))
		for k, v := range b.TfIdf {
			tfidfSlice = append(tfidfSlice, mapEntry{k, v})
		}
		sort.Slice(tfidfSlice, func(i, j int) bool {
			return tfidfSlice[i].value > tfidfSlice[j].value
		})

		xvals := make([]string, 0)
		for i := 0; i < int(float32(len(b.TfIdf))*float32(RelatibilityMod)); i++ {
			xvals = append(xvals, tfidfSlice[i].key)
		}

		b.TfIdf = nil
		b.XValTerms = xvals
		tempBizMap[b.BusinessID] = b
	}
	Businesses = tempBizMap

	log.Println("TF-IDF of Business Review Terms sorted")

}

// Iterates trough all of the businesses then iterates through all of it's td-idf for terms,
// it returns the x most valuable terms (respective of the total term keys in the map)
// and adds this to a Global Relatibility key map
func addMostRelevantTermsKeyMap() {
	tempKeyMap := NewHashMap()
	for _, b := range Businesses {
		tempArr := b.XValTerms
		for i := range tempArr {
			//tempKeyMap[k] = append(tempKeyMap[k], b.BusinessID)
			// Return X most valuable tf-idf to use in global relatability index
			tempKeyMap.Add(tempArr[i], b.BusinessID)
			break
		}

	}

	log.Printf("Most relevant terms added to key map, length: %d", tempKeyMap.size)

	TermKeyMap = tempKeyMap
}

// Gets list of random n numbers of businesses and returns them as a list
// Called by HTTP handler for returning a random list of businesses to select from
func getRandomBusinessList(n int) []bizTuple {
	randomBizList := make([]bizTuple, 0, n)
	for i := 0; i < n; i++ {

		r := rand.Intn(len(Businesses) - 1)
		keyMapIdx := BusinessKeyMap[r]
		randomBizList = append(randomBizList, bizTuple{Businesses[keyMapIdx].Name, Businesses[keyMapIdx].BusinessID})
	}
	log.Printf("Random businesses: %v\n", randomBizList)
	return randomBizList
}

// Returns 2 random relatable businesses
func findRelatableBusinesses(businessID string) []Business {
	log.Println("Finding relatable businesses...")
	relatableKeys := Businesses[businessID].XValTerms // get relatable terms for this business
	relatableBusinesses := make([]Business, 0, 2)     // instantiate 2 businesses to return

	bid1, bid2 := "k", "k"
	found := false
	for !found {
		key1, key2 := "k", "k" // starter keys
		for key1 == key2 {
			r1 := rand.Intn(len(relatableKeys) - 1)
			r2 := rand.Intn(len(relatableKeys) - 1)
			key1 = relatableKeys[r1]
			key2 = relatableKeys[r2]
		}

		bid1, _ = TermKeyMap.Get(key1)
		bid2, _ = TermKeyMap.Get(key2)
		if bid1 != "" && bid2 != "" && bid1 != bid2 && bid1 != businessID && bid2 != businessID {
			found = true
		}
	}

	relatableBusinesses = append(relatableBusinesses, Businesses[bid1])
	relatableBusinesses = append(relatableBusinesses, Businesses[bid2])
	return relatableBusinesses
}
