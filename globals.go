package main

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
	TermFrequency    map[string]float32 `json:"term_frequency"`
	TfIdf            map[string]float32 `json:"tf_idf"`
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
	ReviewTotal           = 25000
	BusinessKeyMap        []string
	ReviewCount           = 100
	RelatibilityMod       = 0.10
	TermKeyMap            = make(map[string][]string)
)

const (
	serviceName        = "Yelp Similarity Web App"
	serviceDescription = "Yep Similarity Web App powered by probabilistic data structures"
	businessPath       = "./json/yelp_academic_dataset_business.json"
	reviewJsonPath     = "./json/yelp_academic_dataset_review.json"
)
