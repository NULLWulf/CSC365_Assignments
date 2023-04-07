package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

// BusinessDataPoint represents a data point in the k-medoids clustering algorithm
type BusinessDataPoint struct {
	BusinessID       string                `json:"business_id"`
	Latitude         float32               `json:"latitude"`
	Longitude        float32               `json:"longitude"`
	ReviewScore      float32               `json:"review_score"`
	FileIndex        int                   `json:"file_index"`
	ClosestNeighbors [4]*BusinessDataPoint `json:"closest_neighbors"`
	Categories       []string              `json:"category"`
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
			dps = append(dps, BusinessDataPoint{BusinessID: b.BusinessID,
				Latitude:    b.Latitude,
				Longitude:   b.Longtitude,
				ReviewScore: b.Stars,
				FileIndex:   l,
				Categories:  b.CategoriesArr})
		}
	}
	log.Printf("BusinessDataPoints loaded: %d", len(dps))
	if k.CCount == 0 {
		k.CCount = 10
	}

	k.PopClusters(dps, k.CCount)
	log.Printf("KmediodsDS built from PSD")
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

// distance - As of assignment uses Haversine distance, but for other data use Euclidean distance
func distance(p1, p2 BusinessDataPoint) float32 {
	// Euclidean distance between two points
	//return float32(math.Sqrt(math.Pow(float64(p1.Latitude-p2.Latitude), 2) + math.Pow(float64(p1.Longitude-p2.Longitude), 2) + math.Pow(float64(p1.ReviewScore-p2.ReviewScore), 2)))
	return haversineDistance(p1, p2)
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
		if isEqualBusinessDataPoint(m, medoid) {
			return i
		}
	}
	return -1
}

func equal(medoids1, medoids2 []BusinessDataPoint) bool {
	// Check if two lists of medoids are equal
	for i, m := range medoids1 {
		if !isEqualBusinessDataPoint(m, medoids2[i]) {
			return false
		}
	}
	return true
}

func isEqualBusinessDataPoint(p1, p2 BusinessDataPoint) bool {
	if p1.BusinessID != p2.BusinessID || p1.Latitude != p2.Latitude || p1.Longitude != p2.Longitude || p1.ReviewScore != p2.ReviewScore || p1.FileIndex != p2.FileIndex {
		return false
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
	medoids := make([]BusinessDataPoint, k)
	clusters := make([]Cluster, k)

	firstIteration := true

	for {
		if firstIteration {
			// Initialize medoids randomly on the first iteration
			for i := 0; i < k; i++ {
				medoids[i] = data[rand.Intn(len(data))]
			}
			firstIteration = false
		} else {
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
			}
		}

		// Store old medoids for convergence checking
		oldMedoids := make([]BusinessDataPoint, k)
		copy(oldMedoids, medoids)

		// Reset clusters and assign medoids
		clusters = make([]Cluster, k)
		for i, medoid := range medoids {
			clusters[i].Medoid = medoid
		}

		// Assign each point to the closest medoid
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

		// Check for convergence
		if equal(medoids, oldMedoids) {
			break
		}
	}
	// TODO - Add 4 closest businesses to each data point

	return clusters
}

func (k *KmediodsDS) BuildGraphFromKM() interface{} {

	log.Printf("Building graph from KMediods data structure...")
	// Create a graph
	graph := NewGraph()

	log.Printf("Adding nodes to graph...")
	for _, c := range k.Clusters {
		for _, p := range c.Points {
			graph.AddNode(p.FileIndex, c.Medoid.FileIndex, p)
		}
	}

	// Calculate 4 closest geographical points to each point, and create edges between them
	log.Printf("Adding edges to graph...")
	for i, c := range k.Clusters {
		log.Printf("Adding edges for cluster %d, cluste size %d", i, len(c.Points))
		for _, p := range c.Points {
			// temporay sorted list of elements in cluster
			temp := make([]BusinessDataPoint, len(c.Points))
			copy(temp, c.Points)
			// Sort the list by distance
			sort.Slice(temp, func(i, j int) bool {
				return distance(p, temp[i]) < distance(p, temp[j])
			})
			temp = temp[1:5]
			for _, t := range temp {
				graph.AddEdge(p.FileIndex, t.FileIndex, jaccardSimilarity(&p, &t))
				graph.AddEdge(t.FileIndex, p.FileIndex, jaccardSimilarity(&p, &t))
			}
		}
		if i == 0 {
			log.Printf("First cluster edges added, breaking...")
			break
		}
	}

	// Running Djikstra's algorithm on the graph
	log.Printf("Running Djikstra's algorithm on graph...")
	// Run Djstrika to root on every point
	log.Printf("Root point: %d", k.Clusters[0].Medoid.FileIndex)
	for _, v := range k.Clusters[0].Points {
		log.Printf("Running Djikstra's algorithm on point %d", v.FileIndex)
		keys, distance := graph.DijkstraShortestPath(v.FileIndex, k.Clusters[0].Medoid.FileIndex)
		log.Printf("Djikstra's algorithm finished, keys: %v, distance: %v", keys, distance)
	}

	graph.SaveGraph()
	return nil
}
