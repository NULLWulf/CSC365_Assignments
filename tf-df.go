package main

import (
	"log"
	"math"
	"math/rand"
	"sort"
)

var (
	RelatibilityMod = .10 // The inreases/decreases the top percentage of terms to be considered relevant
)

var TermKeyMap map[string][]string

// RemoveNullReviewsFromBusinesses removes businesses with no reviews from the Businesses array by
// creating a new array and copying over the businesses with reviews
// It also calculates the term frequency for each term in the reviews of the business
// and increments the document frequency for each term
func removeNullReviewsCalculateFrequency() {
	newBusinesses := make(map[string]Business)
	tdf := make(map[string]int)

	for _, b := range Businesses {
		if len(b.ReviewTermsCount) != 0 {
			delete(b.ReviewTermsCount, "") // delete empty string key if it exists
			tCount := 0                    // total numbers of terms in the reviews of
			for _, count := range b.ReviewTermsCount {
				tCount += count
			}
			b.termFrequency = make(map[string]float32)
			for k, v := range b.ReviewTermsCount {
				// Calculate Term Frequency
				b.termFrequency[k] = float32(v) / float32(tCount)
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

func calculatetfIdf() {
	for i, b := range Businesses {
		b.TfIdf = make(map[string]float32)
		for k, v := range b.termFrequency {
			b.TfIdf[k] = v * float32(math.Log(float64(len(Businesses))/float64(TermDocumentFrequency[k])))
		}
		Businesses[i] = b

	}
	log.Printf("Businesses: %d, TermDocumentFrequencyNumbers: %d\n", len(Businesses), len(TermDocumentFrequency))
}

func sortTfIdf() {
	log.Println("Sorting TF-IDF...")

	type mapEntry struct {
		key   string
		value float32
	}

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

func addMostRelevantTermsKeyMap() map[string][]string {
	tempKeyMap := make(map[string][]string)
	// TermKeyMap = make(map[string][]string)
	for _, b := range Businesses {
		i := 0
		for k := range b.TfIdf {
			tempKeyMap[k] = append(tempKeyMap[k], b.BusinessID)
			if i > int(float32(len(b.TfIdf))*float32(RelatibilityMod)) {
				break
			}
			i++
		}
	}
	log.Printf("Most relevant terms added to key map, length: %d", len(tempKeyMap))
	return tempKeyMap
}

// businessId is a key in this instance, keys are terms in the instance
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

func getRandomBusiness(n int) interface{} {
	type bizTuple struct {
		businessName string
		businessID   string
	}
	randomBizList := make([]bizTuple, 0, n)
	for i := 0; i < n; i++ {

		r := rand.Intn(len(Businesses) - 1)
		keyMapIdx := BusinessKeyMap[r]
		randomBizList = append(randomBizList, bizTuple{Businesses[keyMapIdx].Name, Businesses[keyMapIdx].BusinessID})
	}
	log.Printf("Random businesses: %v\n", randomBizList)
	return randomBizList
}

func findRelatableBusinesses(businessID string) interface{} {
	relatableKeys := getRelatableByKey(businessID)   // keys from this business
	relatableBusinesses := make(map[string]Business) // leys to store matchinf relatable businesses by terms

	for found := 0; found == 2; {
		tryKey := relatableKeys[rand.Intn(len(relatableKeys)-1)]
		for _, bID := range TermKeyMap[tryKey] {
			if bID != businessID {
				relatableBusinesses[bID] = Businesses[bID]
				found++
				break
			}
		}
	}

	return relatableBusinesses
}
