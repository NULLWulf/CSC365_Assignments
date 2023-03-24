package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bbalet/stopwords"
)

// Reads Businesses JSON data
func readBusinessesJson() {
	log.Fatal("Deprecated: Use readBusinessesJson2() instead.")
	log.Println("Loading Business JSON data...")
	t := 0
	// Create directory for fileblock if it does not exist
	err := os.MkdirAll("fileblock", 0777)
	if err != nil {
		log.Fatal(err)
	}

	// Bussiness ID list
	var businessIDList []string
	// Read the file containing business information
	file, err := os.ReadFile(businessPath) // Reads entire file to file object
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Businesses Loaded: %d", len(file))
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
				file := fmt.Sprintf("fileblock/%d.json", t)
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
				businessIDList = append(businessIDList, business.BusinessID)

				t++

			}
			// }()

		}
		i = j + 1
		if t == 10 {
			break
		}
	}
	log.Printf("Businesses Loaded: %d", t)
	// businessMap := NewEHT2(10)
	// for _, id := range businessIDList {
	// 	businessMap.insert(id)
	// }
	// log.Println("Businesses Loaded: ", businessMap.DirectorySize)

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

/*
	ReadBusinessJSON2

Read the Business JSON file, get the businesses that are open, have more than 100 reviews, and are restaurants.
Then write the businesses to a JSON file, with a limit of 10000 businesses.  As it reads through the JSON file,
it unmarshals the JSON into a Business struct, and then marshals the struct into a JSON file.
Adds the filename to the Extensible Hash Table, which is used as a file index..  During the process,
it also gets the relevant information from the business struct and stores it in a BusinessDataPoint struct
for the K-medoids algorithm.
*/
func ReadBusinessJSON2() []BusinessDataPoint {
	InstantiateFileBlock()
	log.Println("Loading Business JSON data...")
	t := 0
	file, err := os.ReadFile("yelp_academic_dataset_business.json")
	if err != nil {
		log.Fatal(err)
	}
	eht := NewEHT2(5000)

	// businessIDList := make([]string, 0)
	BusinessDataPoints := make([]BusinessDataPoint, 0)
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
				// Write business to JSON file
				file := fmt.Sprintf("fileblock/%d.json", t)
				businessFile, err := os.Create(file)
				if err != nil {
					log.Fatal(err)
				}
				businessJson, err := json.Marshal(business)
				if err != nil {
					log.Fatal(err)
				}
				_, err = businessFile.Write(businessJson)
				if err != nil {
					log.Fatal(err)
				}
				var b Business
				err = json.Unmarshal(businessJson, &b)
				businessFile.Close()
				BusinessDataPoints = append(BusinessDataPoints, BusinessDataPoint{BusinessID: business.BusinessID, Latitude: b.Latitude, Longitude: b.Longtitude, ReviewScore: b.Stars, FileIndex: t})
				eht.insert(t)
				t++
			}

		}
		i = j + 1
		if t == 10000 {
			break
		}
	}

	// Insert Business IDs into Extensible Hash Table
	log.Printf("Businesses Loaded: %d", t)

	err = eht.saveToDisk("artifacts")
	if err != nil {
		// make custom error
		log.Fatal(errors.New("Failed to save to disk:" + err.Error()))
	}
	log.Printf("EHT saved to disk")
	return BusinessDataPoints
}

func ReadDirectory(url string) {

	fileList, err := os.ReadDir(url)
	if err != nil {
		log.Printf("Failed to read directory")
	}

	for i, e := range fileList {
		log.Printf("File %d : %s", i, e.Name())
	}
}

func DeleteDirectory(url string) {
	err := os.RemoveAll(url)
	if err != nil {
		log.Printf("Failed to delete directory")
	}
}

func InstantiateFileBlock() {
	// See if directory exists and if it does delete it
	log.Printf("Instantiating fileblock directory...")
	if _, err := os.Stat("fileblock"); !os.IsNotExist(err) {
		log.Printf("fileblock directory already exists.  Deleting...")
		DeleteDirectory("fileblock")
	}
	log.Printf("fileblock directory deleted.  Creating new directory...")
	// Create directory for fileblock if it does not exist
	err := os.MkdirAll("fileblock", 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func GetRandomFileNames(dirPath string, numFiles int) ([]string, error) {
	rand.Seed(time.Now().UnixNano()) // set random seed

	// Get list of all files in directory
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	// Shuffle the files randomly
	rand.Shuffle(len(files), func(i, j int) {
		files[i], files[j] = files[j], files[i]
	})

	// Get the names of the first `numFiles` files
	var result []string
	for i := 0; i < numFiles && i < len(files); i++ {
		if !files[i].IsDir() {
			result = append(result, files[i].Name())
		}
	}

	return result, nil
}

func LoadBusinessFromFile(businessID string) Business {
	file, err := os.ReadFile(fmt.Sprintf("fileblock/%s", businessID))
	if err != nil {
		log.Fatal(err)
	}
	var business Business
	err = json.Unmarshal(file, &business)
	if err != nil {
		log.Fatal(err)
	}
	return business
}