package main

import (
	"encoding/gob"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type ExtensibleHashTable struct {
	BucketSize    int       `json:"bucket_size"`
	BucketArr     []*Bucket `json:"bucket_arr"`
	DirectorySize int       `json:"directory_size"`
	GlobalDepth   int       `json:"global_depth"`
}

// NewEHT2 EHT constructor specifying an EHT with a specific bucketsize
func NewEHT2(bucketSize int) *ExtensibleHashTable {
	gd := 1
	bs := bucketSize
	ds := int(math.Pow(2, float64(gd)))
	bucketArr := make([]*Bucket, ds)
	for i := 0; i < ds; i++ {
		bucketArr[i] = NewBucket2(bs)
	}

	return &ExtensibleHashTable{
		BucketSize:    bs,
		BucketArr:     bucketArr,
		DirectorySize: ds,
		GlobalDepth:   gd,
	}
}

// NewEHT constructor for instantiating an EHT without a specific bucket size
func NewEHT() *ExtensibleHashTable {
	gd := 1
	bs := 4
	ds := int(math.Pow(2, float64(gd)))
	bucketArr := make([]*Bucket, ds)
	for i := 0; i < ds; i++ {
		bucketArr[i] = NewBucket()
	}

	return &ExtensibleHashTable{
		BucketSize:    bs,
		BucketArr:     bucketArr,
		DirectorySize: ds,
		GlobalDepth:   gd,
	}
}

// getSize method for getting number of buckets
func (eht *ExtensibleHashTable) getSize() int {
	return eht.BucketSize
}

// getGD method for getting global depth of EHT
func (eht *ExtensibleHashTable) getGD() int {
	return eht.GlobalDepth
}

// getDS method for getting the directory size of EHT
func (eht *ExtensibleHashTable) getDS() int {
	return eht.DirectorySize
}

// method for performing a split when a bucket is full
func (eht *ExtensibleHashTable) split(key int) {
	index := eht.FNVHash(key)                   // Get index of bucket to split
	oldBucket := eht.BucketArr[index]           // Get bucket to split
	newBucket := NewBucket2(eht.BucketSize)     // Create new bucket
	newBucket.LocalDepth = oldBucket.LocalDepth // Set local depth of new bucket to local depth of old bucket

	// Double the size of the directory if the global depth is equal to the local depth
	if eht.GlobalDepth == oldBucket.LocalDepth {
		eht.doubleArray()
	}

	newIndex := eht.FNVHash(key) // Get index of new bucket
	toBeMoved := []int{}         // Array of values to be moved to new bucket

	// Move values to new bucket
	for i := 0; i < eht.BucketSize; i++ {
		if eht.FNVHash(oldBucket.ValueArr[i]) == newIndex {
			toBeMoved = append(toBeMoved, oldBucket.ValueArr[i])
			oldBucket.ValueArr[i] = 0
		}
	}

	// Insert values into new bucket
	// Remove values from old bucket
	for _, val := range toBeMoved {
		newBucket.insert(val)
		oldBucket.remove(val)
	}

	newBucket.LocalDepth++
	eht.BucketArr[newIndex] = newBucket

	var indexSameConnection int
	// Update local depth of bucket with same connection
	if newIndex >= eht.DirectorySize/2 {
		indexSameConnection = newIndex - int(math.Pow(2, float64(eht.GlobalDepth-1)))
	} else {
		indexSameConnection = newIndex + int(math.Pow(2, float64(eht.GlobalDepth-1)))
	}
	eht.BucketArr[indexSameConnection].LocalDepth++
}

// doubleArray method for doubling the size of the directory
// when the global depth is equal to the local depth
func (eht *ExtensibleHashTable) doubleArray() {
	newDirectorySize := eht.DirectorySize * 2
	newArray := make([]*Bucket, newDirectorySize)

	// Update local depth of existing buckets
	for i := 0; i < eht.DirectorySize; i++ {
		eht.BucketArr[i].LocalDepth++
	}

	// Rehash all items in the hash table
	for i := 0; i < eht.DirectorySize; i++ {
		for j := 0; j < eht.BucketSize; j++ {
			if eht.BucketArr[i].ValueArr[j] != 0 {
				index := eht.FNVHash(eht.BucketArr[i].ValueArr[j])
				newArray[index] = eht.BucketArr[i]
			}
		}
	}

	eht.BucketArr = newArray
	eht.GlobalDepth++
	eht.DirectorySize = newDirectorySize
}

// insert method for inserting a value into the EHT
func (eht *ExtensibleHashTable) insert(key int) {
	index := eht.FNVHash(key)
	// If bucket is full, split the bucket
	if eht.BucketArr[index].isFull() {
		eht.split(key)
		eht.insert(key)
	} else { // otherwise, insert the value into the bucket
		eht.BucketArr[index].insert(key)
	}
}

// find method for finding a value in the EHT
func (eht *ExtensibleHashTable) find(key int) bool {
	index := eht.FNVHash(key)
	if eht.BucketArr[index].find(key) {
		return true
	}
	return false
}

// remove method for removing a value from the EHT
func (eht *ExtensibleHashTable) remove(key int) bool {
	index := eht.FNVHash(key)
	if eht.BucketArr[index].find(key) {
		eht.BucketArr[index].remove(key)
		return true
	}
	return false
}

// FNVHash  Originally was using Go build FNVO hashing function but it seemed
// to cause some performance issues that were beyond the scope this assignment
func (eht *ExtensibleHashTable) FNVHash(key int) int {
	return (key * 16777619) * 37 % eht.DirectorySize
}

// print Utilit helper function for displaying EHT metadata
func (eht *ExtensibleHashTable) print() {
	dis := strings.Builder{}
	dis.WriteString("Directory Size: " + strconv.Itoa(eht.DirectorySize) + "\n")
	dis.WriteString("Global Depth: " + strconv.Itoa(eht.GlobalDepth) + "\n")
	dis.WriteString("Bucket Size: " + strconv.Itoa(eht.BucketSize) + "\n")
	log.Printf(dis.String())
}

// saveToDisk Save the EHT to disc as a binary gob
func (eht *ExtensibleHashTable) saveToDisk(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(eht)
	if err != nil {
		return err
	}

	return nil
}

// deserialize deserializes the ExtensibleHashTable from the given file path.
func deserialize(filePath string) (*ExtensibleHashTable, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var eht ExtensibleHashTable
	err = decoder.Decode(&eht)
	if err != nil {
		return nil, err
	}

	return &eht, nil
}
