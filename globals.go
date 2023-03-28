package main

/*
 * This file contains all the global variables and structs used in the program
 * But is not all inclusive of the program.
 */

type Business struct {
	BusinessID       string             `json:"business_id"`
	Name             string             `json:"name"`
	City             string             `json:"city"`
	State            string             `json:"state"`
	Stars            float32            `json:"stars"`
	ReviewCount      int                `json:"review_count"`
	IsOpen           int                `json:"is_open"`
	Categories       string             `json:"categories"`
	CategoriesArr    []string           `json:"-"`
	Latitude         float32            `json:"latitude"`
	Longtitude       float32            `json:"longitude"`
	ReviewTermsCount map[string]int     `json:"-"`
	TermCountTotal   int                `json:"-"`
	TermFrequency    map[string]float32 `json:"-"`
	TfIdf            map[string]float32 `json:"-"`
	XValTerms        []string           `json:"-"`
	FileId           uint64             `json:"file_id"`
}

type Review struct {
	ReviewID   string  `json:"review_id"`
	UserID     string  `json:"user_id"`
	BusinessID string  `json:"business_id"`
	Stars      float32 `json:"stars"`
	Text       string  `json:"text"`
}

type mapEntry struct {
	key   string
	value float32
}

type program struct {
	Name        string
	DisplayName string
	Description string
}

type bizTuple struct {
	BusinessName string `json:"business_name"`
	BusinessID   string `json:"business_id"`
}

var (
	Businesses            = make(map[string]Business)
	TermDocumentFrequency = make(map[string]int)
	ReviewTotal           = 50000
	BusinessKeyMap        []string
	ReviewCount           = 50
	RelatibilityMod       = 0.25 //TermKeyMap            = make(map[string][]string)
	TermKeyMap            = NewHashMap()
)

const (
	serviceName        = "Yelp Similarity Web App"
	serviceDescription = "Yep Similarity Web App powered by probabilistic data structures"
	businessPath       = "./json/yelp_academic_dataset_business.json"
	reviewJsonPath     = "./json/yelp_academic_dataset_review.json"
)
