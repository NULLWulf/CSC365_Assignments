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

func NewEHT() *ExtensibleHashTable {
	gd := 1
	bs := 4
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

func (eht *ExtensibleHashTable) getSize() int {
	return eht.BucketSize
}

func (eht *ExtensibleHashTable) getGD() int {
	return eht.GlobalDepth
}

func (eht *ExtensibleHashTable) getDS() int {
	return eht.DirectorySize
}

func (eht *ExtensibleHashTable) split(key string) {
	index := eht.FNVHash(key)
	oldBucket := eht.BucketArr[index]
	newBucket := NewBucket2(eht.BucketSize)
	newBucket.LocalDepth = oldBucket.LocalDepth

	if eht.GlobalDepth == oldBucket.LocalDepth {
		eht.doubleArray()
	}

	newIndex := eht.FNVHash(key)
	toBeMoved := []string{}

	for i := 0; i < eht.BucketSize; i++ {
		if eht.FNVHash(oldBucket.ValueArr[i]) == newIndex {
			toBeMoved = append(toBeMoved, oldBucket.ValueArr[i])
			oldBucket.ValueArr[i] = ""
		}
	}

	for _, val := range toBeMoved {
		newBucket.insert(val)
		oldBucket.remove(val)
	}

	newBucket.LocalDepth++
	eht.BucketArr[newIndex] = newBucket

	var indexSameConnection int
	if newIndex >= eht.DirectorySize/2 {
		indexSameConnection = newIndex - int(math.Pow(2, float64(eht.GlobalDepth-1)))
	} else {
		indexSameConnection = newIndex + int(math.Pow(2, float64(eht.GlobalDepth-1)))
	}
	eht.BucketArr[indexSameConnection].LocalDepth++
}

func (eht *ExtensibleHashTable) doubleArray() {
	newDirectorySize := eht.DirectorySize * 2
	newArray := make([]*Bucket, newDirectorySize)

	for i := 0; i < eht.DirectorySize; i++ {
		newArray[i] = eht.BucketArr[i]
		newArray[i+eht.DirectorySize] = eht.BucketArr[i]
	}

	eht.BucketArr = newArray
	eht.GlobalDepth++
	eht.DirectorySize = newDirectorySize
}

func (eht *ExtensibleHashTable) insert(key string) {
	index := eht.FNVHash(key)
	if eht.BucketArr[index].isFull() {
		eht.split(key)
		eht.insert(key)
	} else {
		eht.BucketArr[index].insert(key)
	}
}

func (eht *ExtensibleHashTable) find(key string) bool {
	index := eht.FNVHash(key)
	if eht.BucketArr[index].find(key) {
		return true
	}
	return false
}

func (eht *ExtensibleHashTable) remove(key string) bool {
	index := eht.FNVHash(key)
	if eht.BucketArr[index].find(key) {
		eht.BucketArr[index].remove(key)
		return true
	}
	return false
}

// FNVOHash
func (eht *ExtensibleHashTable) FNVHash(key string) int {
	hash := uint32(5381)
	for _, char := range key {
		hash = ((hash << 5) + hash) + uint32(char) // hash * 33 + char
	}
	return int(hash % uint32(eht.DirectorySize))
}

func (eht *ExtensibleHashTable) print() {
	dis := strings.Builder{}
	dis.WriteString("Directory Size: " + strconv.Itoa(eht.DirectorySize) + "\n")
	dis.WriteString("Global Depth: " + strconv.Itoa(eht.GlobalDepth) + "\n")
	dis.WriteString("Bucket Size: " + strconv.Itoa(eht.BucketSize) + "\n")
	log.Printf(dis.String())
}

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
