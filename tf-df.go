package main

import (
	"log"
	"math"
	"sort"
)

// RemoveNullReviewsFromBusinesses removes businesses with no reviews from the Businesses array by
// creating a new array and copying over the businesses with reviews
// It also calculates the term frequency for each term in the reviews of the business
// and increments the document frequency for each term
func removeNullReviewsCalculateFrequency() {
	var newBusinesses []Business
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
			newBusinesses = append(newBusinesses, b)
		}
	}
	// TermDocumentFrequency = make(map[string]int)
	TermDocumentFrequency = tdf
	Businesses = newBusinesses
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

	// Sort the TfIdf maps in each Business in descending order of value
	for i, b := range Businesses {
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
		Businesses[i].TfIdf = sortedMap
	}

	log.Println("TF-IDF of Business Review Terms sorted")

}
