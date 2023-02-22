package main

import (
	"log"
	"math"
	"math/rand"
	"sort"
)

// Make a map of strings to floats
//

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
			// Replace b.TermFrequency with actual hashmap
			b.TermFrequency = make(map[string]float32) // <-------------
			for k, v := range b.ReviewTermsCount {
				// Calculate Term Frequency
				//
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
		sortedMap := make(map[string]float32, len(tfidfSlice))
		for _, entry := range tfidfSlice {
			sortedMap[entry.key] = entry.value
		}
		b.TfIdf = sortedMap
		tempBizMap[b.BusinessID] = b
	}
	Businesses = tempBizMap

	log.Println("TF-IDF of Business Review Terms sorted")

}

// Iterates trough all of the businesses then iterates through all of it's td-idf for terms,
// it returns the x most valuable terms (respective of the total term keys in the map)
// and adds this to a Global Relatibility key map
func addMostRelevantTermsKeyMap() {
	TermKeyMap = make(map[string][]string)
	tempKeyMap := make(map[string][]string)
	for _, b := range Businesses {
		i := 0
		for k := range b.TfIdf {
			tempKeyMap[k] = append(tempKeyMap[k], b.BusinessID)
			// Return X most valuable tf-idf to use in global relatability index
			if i > int(float32(len(b.TfIdf))*float32(RelatibilityMod)) {
				break
			}
			i++
		}
	}
	log.Printf("Most relevant terms added to key map, length: %d", len(tempKeyMap))
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
	return randomBizList
}

// Called from an http handler, this finds the relatable business by gettings
// it's most valuable td-idf term-keys into a list, then it randomly selects
// an element from said list to get a specifc term-key.  Then it finds the term
// as a key in the TermKeysMap, in which the value at the map element is an array of
// business Ids whom share that term for relatively high tf-idf.
// It then iterates through and picks a business in tthe TermKeyMap subset.
func findRelatableBusinesses(businessID string) []Business {
	log.Println("Finding relatable businesses...")
	relatableKeys := getRelatableByKey(businessID) // keys from this business
	relatableBusinesses := make([]Business, 0, 2)  // 2 businesses to return

	for found := 0; found < 2; {
		tryKey := relatableKeys[rand.Intn(len(relatableKeys)-1)]
		for _, bID := range TermKeyMap[tryKey] { // term key not set as a global variable at this point
			if bID != businessID {
				// check if key is already in relatableBusinesses
				keyExists := false
				for _, b := range relatableBusinesses {
					if b.BusinessID == bID {
						keyExists = true
						break
					}
				}
				// add business to relatableBusinesses if key is not already in it
				if !keyExists {
					relatableBusinesses = append(relatableBusinesses, Businesses[bID])
					found++
					break
				}
			}
		}
	}

	return relatableBusinesses
}

// Returns X most valuable tf-idf terms for a given business
// By default it will get the most 10% valuable tf-idfs however this variable can
// tweaked to get a larger range
func getRelatableByKey(key string) []string {

	tBussiness := Businesses[key]
	tTDIDF := tBussiness.TfIdf
	i := 0
	tKeys := make([]string, 0, len(tTDIDF))
	for k := range tTDIDF {
		if i > int(float32(len(tTDIDF))*float32(RelatibilityMod)) {
			break
		}
		tKeys = append(tKeys, k)
		i++
	}
	// Get the top 10% of terms in the business
	return tKeys
}
