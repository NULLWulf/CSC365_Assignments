package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
)

// BusinessDataPoint represents a data point in the k-medoids clustering algorithm
type BusinessDataPoint struct {
	BusinessID  string  `json:"business_id"`
	Latitude    float32 `json:"latitude"`
	Longitude   float32 `json:"longitude"`
	ReviewScore float32 `json:"review_score"`
	FileIndex   int     `json:"file_index"`
}

type Cluster struct {
	Medoid BusinessDataPoint
	Points []BusinessDataPoint
	ID     int
}

type KmediodsDS struct {
	Clusters []Cluster
	CCount   int
}

type SimClusterResponse struct {
	ClusterSize int      `json:"cluster_size"`
	Selected    Business `json:"business_selected"`
	Medoid      Business `json:"business_medoid"`
	Similar     Business `json:"businesses_similar"`
}

// BuildFromPSD builds a KmediodsDS from ExtensibleHashTable which services a file index in this case
func (k *KmediodsDS) BuildFromPSD() {
	k.SetIfNull()
	log.Printf("Building KmediodsDS from PSD")
	log.Printf("Loading EHT from disk...")
	eht, err := deserialize("artifacts") //
	if err != nil {
		log.Fatal(err)
	}
	dps := make([]BusinessDataPoint, 0)
	log.Printf("Populating BusinessDataPoints from EHT...")
	seen := make(map[int]bool) // keeps track of seen FileIndex values
	// iterate over all buckets and values in the EHT and populate BusinessDataPoints
	for _, v := range eht.BucketArr {
		for _, l := range v.ValueArr {
			// convert number to actual string
			a := strconv.Itoa(l)
			// load Business from file
			b := LoadBusinessFromFile(a)
			if seen[l] {
				continue // skip to next iteration of inner loop
			}
			seen[l] = true // mark FileIndex as seen
			dps = append(dps, BusinessDataPoint{BusinessID: b.BusinessID, Latitude: b.Latitude, Longitude: b.Longtitude, ReviewScore: b.Stars, FileIndex: l})
		}
	}
	log.Printf("BusinessDataPoints loaded: %d", len(dps))
	if k.CCount == 0 {
		k.CCount = 10
	}

	k.PopClusters(dps, k.CCount)
	log.Printf("KmediodsDS built from PSD")
}

// KMedoids performs k-medoids clustering on the given data set
// and returns the resulting clusters
func KMedoids(data []BusinessDataPoint, k int) []Cluster {
	// Initialize medoids randomly
	medoids := make([]BusinessDataPoint, k)
	// get random k values to service as initial comparison points
	for i := 0; i < k; i++ {
		medoids[i] = data[rand.Intn(len(data))]
	}

	// Assign each point to the closest medoid
	clusters := make([]Cluster, k)
	// Set the temporary cluster medoids to the randomly selected medoids
	for i, medoid := range medoids {
		clusters[i].Medoid = medoid
	}
	// Iterate over all data points and assign them to the closest medoid
	// by computing the distance between the point and each medoid
	// and selecting the medoid with the lowest distance
	for _, point := range data {
		minDist := math.MaxFloat32
		var closestMedoid BusinessDataPoint
		for _, medoid := range medoids {
			dist := distance(point, medoid)
			if float64(dist) < minDist {
				minDist = float64(dist)
				closestMedoid = medoid
			}
		}
		clusterIndex := findIndex(medoids, closestMedoid)
		clusters[clusterIndex].Points = append(clusters[clusterIndex].Points, point)
	}

	// Update medoids by computing the cost of each point in each cluster
	// and selecting the point with the lowest cost as the new medoid
	for i := 0; i < k; i++ {
		minCost := math.MaxFloat32
		var newMedoid BusinessDataPoint
		for _, point := range clusters[i].Points {
			cost := computeCost(clusters[i].Points, point)
			if float64(cost) < minCost {
				minCost = float64(cost)
				newMedoid = point
			}
		}
		medoids[i] = newMedoid
		clusters[i].Medoid = newMedoid
	}

	// Repeat until convergence
	for {
		oldMedoids := make([]BusinessDataPoint, k)
		copy(oldMedoids, medoids)

		// Assign each point to the closest medoid
		clusters = make([]Cluster, k)
		// Set the temporary cluster medoids to the randomly selected medoids
		for i, medoid := range medoids {
			clusters[i].Medoid = medoid
		}
		for _, point := range data {
			minDist := math.MaxFloat32
			var closestMedoid BusinessDataPoint
			for _, medoid := range medoids {
				dist := distance(point, medoid)
				if float64(dist) < minDist {
					minDist = float64(dist)
					closestMedoid = medoid
				}
			}
			clusterIndex := findIndex(medoids, closestMedoid)
			clusters[clusterIndex].Points = append(clusters[clusterIndex].Points, point)
		}

		// Update medoids by computing the cost of each point in each cluster
		// and selecting the point with the lowest cost as the new medoid
		for i := 0; i < k; i++ {
			minCost := math.MaxFloat32
			var newMedoid BusinessDataPoint
			for _, point := range clusters[i].Points {
				cost := computeCost(clusters[i].Points, point)
				if float64(cost) < minCost {
					minCost = float64(cost)
					newMedoid = point
				}
			}
			medoids[i] = newMedoid
			clusters[i].Medoid = newMedoid
		}

		// Check for convergence
		if equal(medoids, oldMedoids) {
			break
		}
	}

	return clusters
}

// GetRandomDataPoints Gets a random cluster and value within said cluster and returns the random list
// of points
func (k *KmediodsDS) GetRandomDataPoints(ct int) []BusinessDataPoint {
	randPoints := make([]BusinessDataPoint, ct)
	seen := make(map[string]bool)
	for i := 0; i < ct; {
		randCluster := k.Clusters[rand.Intn(len(k.Clusters))]
		randPoint := randCluster.Points[rand.Intn(len(randCluster.Points))]
		coord := fmt.Sprintf("%v,%v", randCluster, randPoint)
		// if we have not seen this point before, add it to the list
		if _, ok := seen[coord]; !ok {
			seen[coord] = true
			randPoints[i] = randPoint
			i++
		}
	}

	return randPoints
}

// FindSimilarBuildResponse finds a similar business to the one passed in
func (k *KmediodsDS) FindSimilarBuildResponse(fileId string) ([]byte, error) {
	var selectedBiz Business
	file, err := os.ReadFile("fileblock/" + fileId + ".json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &selectedBiz)
	log.Printf("fileblock/%s.json: %+v", fileId, selectedBiz)
	// convert business to data point
	bdps := BusinessDataPoint{
		BusinessID:  selectedBiz.BusinessID,
		Latitude:    selectedBiz.Latitude,
		Longitude:   selectedBiz.Longtitude, // corrected typo: Longtitude -> Longitude
		ReviewScore: selectedBiz.Stars,
	}

	log.Print("Finding cluster for business: ", bdps.BusinessID)
	// set initial distance to max float32, used to hold the smallest distance
	// between the business and the medoid of a cluster
	var minDst float32 = math.MaxFloat32
	// sets the cluster that is the most similar to the business
	var simCluster Cluster
	// iterate through the clusters, looking at the cluster medoid
	// to determine which cluster is the most similar to the business
	for ki, v := range k.Clusters {
		m := v.Medoid
		eucDst := distance(bdps, m)
		if eucDst < minDst {
			minDst = eucDst
			simCluster = k.Clusters[ki]
		}
	}
	log.Printf("Found similar cluster: %+v", simCluster.ID)

	var minDstInCluster float32 = math.MaxFloat32
	var simBiz Business

	log.Printf("Finding business in cluster: %+v", simCluster.ID)
	biza := simCluster.Points[0]
	for _, bizl := range simCluster.Points {
		dst := distance(bdps, bizl)
		// if dst is less than minDist and the business is not the same as the selected business
		// set the minDist to the new distance and set the similar business to the new business
		if dst < minDstInCluster && bizl.BusinessID != bdps.BusinessID {
			minDstInCluster = dst
			biza = bizl
		}
	}

	// load the similar business from the file
	simBiz = LoadBusinessFromFile(strconv.Itoa(biza.FileIndex))
	log.Printf("Found similar business with Euc distnace of : %f", minDstInCluster)
	log.Printf("Found similar business: %+v", simBiz)

	// Compare Lat Long Stars of 2 Businesses
	log.Printf("Comparing Lat Long Stars of 2 Businesses")
	log.Printf("Selected Business: %+v", selectedBiz)
	log.Printf("Similar Business: %+v", simBiz)

	// Builds a response for the endpoint based on the similar business and the selected business
	b, err := BuildRelatedClusterResponse(&selectedBiz, &simBiz, &simCluster)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// BuildRelatedClusterResponse builds the response for the related cluster calculation and endpoint
func BuildRelatedClusterResponse(sim *Business, sel *Business, clus *Cluster) ([]byte, error) {

	med := clus.Medoid                                          // get the medoid business for the cluster
	medBiz := LoadBusinessFromFile(strconv.Itoa(med.FileIndex)) // Load the medoid business from the file

	resp := SimClusterResponse{
		Similar:     *sim,
		Selected:    *sel,
		Medoid:      medBiz,
		ClusterSize: len(clus.Points),
	}

	// Marshal to bytes
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func distance(p1, p2 BusinessDataPoint) float32 {
	// Euclidean distance between two points
	return float32(math.Sqrt(math.Pow(float64(p1.Latitude-p2.Latitude), 2) + math.Pow(float64(p1.Longitude-p2.Longitude), 2) + math.Pow(float64(p1.ReviewScore-p2.ReviewScore), 2)))
}

func computeCost(cluster []BusinessDataPoint, point BusinessDataPoint) float32 {
	// Sum of distances between a point and all other points in the cluster
	var cost float32
	for _, p := range cluster {
		cost += distance(point, p)
	}
	return cost
}

func findIndex(medoids []BusinessDataPoint, medoid BusinessDataPoint) int {
	// Find the index of a medoid in the list of medoids
	for i, m := range medoids {
		if m == medoid {
			return i
		}
	}
	return -1
}

func equal(medoids1, medoids2 []BusinessDataPoint) bool {
	// Check if two lists of medoids are equal
	for i, m := range medoids1 {
		if m != medoids2[i] {
			return false
		}
	}
	return true
}

func (k *KmediodsDS) saveKMDStoDisc(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(&k)
	if err != nil {
		return err
	}

	return nil
}

func (k *KmediodsDS) loadKMDStoDisc(filePath string) error {
	if k == nil {
		k = &KmediodsDS{}
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&k)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return nil
}

func (k *KmediodsDS) SetClusterCt(ct int) {
	k.CCount = ct
}

func (k *KmediodsDS) SetIfNull() {
	if k == nil {
		k = &KmediodsDS{}
	}
}

func (k *KmediodsDS) PopClusters(data []BusinessDataPoint, l int) {
	k.Clusters = make([]Cluster, l)
	k.Clusters = KMedoids(data, l)
}

func KMedoids(data []BusinessDataPoint, k int) []Cluster {
	// Initialize medoids randomly
	medoids := make([]BusinessDataPoint, k)
	// get random k values to service as initial comparison points
	for i := 0; i < k; i++ {
		medoids[i] = data[rand.Intn(len(data))]
	}

	// Initialize clusters
	clusters := make([]Cluster, k)
	for i, medoid := range medoids {
		clusters[i].Medoid = medoid
	}

	// Initialize variables for convergence check
	oldMedoids := make([]BusinessDataPoint, k)
	firstIteration := true

	// Repeat until convergence
	for {
		// Assign each point to the closest medoid and update medoids
		for i, point := range data {
			minDist := math.MaxFloat32
			var closestMedoid BusinessDataPoint
			var clusterIndex int
			for j, medoid := range medoids {
				dist := distance(point, medoid)
				if float64(dist) < minDist {
					minDist = float64(dist)
					closestMedoid = medoid
					clusterIndex = j
				}
			}
			clusters[clusterIndex].Points = append(clusters[clusterIndex].Points, point)

			// Update medoid for this cluster
			minCost := math.MaxFloat32
			var newMedoid BusinessDataPoint
			for _, p := range clusters[clusterIndex].Points {
				cost := computeCost(clusters[clusterIndex].Points, p)
				if float64(cost) < minCost {
					minCost = float64(cost)
					newMedoid = p
				}
			}
			medoids[clusterIndex] = newMedoid
			clusters[clusterIndex].Medoid = newMedoid
		}

		// Check for convergence
		if !firstIteration && equal(medoids, oldMedoids) {
			break
		}
		copy(oldMedoids, medoids)

		// Reset cluster points for next iteration
		for i := range clusters {
			clusters[i].Points = make([]BusinessDataPoint, 0)
		}

		// Update flag and loop variables
		firstIteration = false
	}

	return clusters
}

