package main

import (
	"encoding/gob"
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

func (k *KmediodsDS) BuildFromPSD() {
	k.SetIfNull()
	log.Printf("Building KmediodsDS from PSD")
	log.Printf("Loading EHT from disk...")
	eht, err := deserialize("artifacts")
	if err != nil {
		log.Fatal(err)
	}
	dps := make([]BusinessDataPoint, 0)
	log.Printf("Populating BusinessDataPoints from EHT...")
	seen := make(map[int]bool) // keep track of seen FileIndex values
	for _, v := range eht.BucketArr {
		for _, l := range v.ValueArr {
			// convert number to actual string
			a := strconv.Itoa(l)
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
	for i := 0; i < k; i++ {
		medoids[i] = data[rand.Intn(len(data))]
	}

	// Assign each point to the closest medoid
	clusters := make([]Cluster, k)
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

	// Repeat until convergence
	for {
		oldMedoids := make([]BusinessDataPoint, k)
		copy(oldMedoids, medoids)

		// Assign each point to the closest medoid
		clusters = make([]Cluster, k)
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

// GetRandomDataPoints Gets a random set of points from the clusters
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
