package main

import "log"

type Bucket struct {
	ValueArr   []int
	LocalDepth int
	Size       int
}

// NewBucket2  - Constructor for Bucket when size is known
func NewBucket2(bucketSize int) *Bucket {
	ld := 1
	sz := bucketSize
	valueArr := make([]int, sz)
	for i := 0; i < sz; i++ {
		valueArr[i] = 0
	}
	return &Bucket{
		ValueArr:   valueArr,
		LocalDepth: ld,
		Size:       sz,
	}
}

// NewBucket - Constructor for Bucket
func NewBucket() *Bucket {
	ld := 1
	sz := 4
	valueArr := make([]int, sz)
	for i := 0; i < sz; i++ {
		valueArr[i] = 0
	}
	return &Bucket{
		ValueArr:   valueArr,
		LocalDepth: ld,
		Size:       sz,
	}
}

// getLD Method for getting local depth of a bucket
func (b *Bucket) getLD() int {
	return b.LocalDepth
}

// getLD Method for getting local depth of a bucket
func (b *Bucket) getSize() int {
	return b.Size
}

// getLD Method for getting local depth of a bucket
func (b *Bucket) getArr() []int {
	return b.ValueArr
}

// insert is a Method for inserting a value in a bucket
func (b *Bucket) insert(value int) {
	for i := 0; i < b.Size; i++ {
		if b.ValueArr[i] == 0 {
			b.ValueArr[i] = value
			break
		}
	}
}

// remove is a method for removing a value in a bucket
func (b *Bucket) remove(value int) {
	for i := 0; i < b.Size; i++ {
		if b.ValueArr[i] == value {
			b.ValueArr[i] = 0
			break
		}
	}
}

// find is a Method for indicating presence of a value in a bucket
func (b *Bucket) find(value int) bool {
	for i := 0; i < b.Size; i++ {
		if b.ValueArr[i] == value {
			return true
		}
	}
	return false
}

// search is a helper Method for debugging
func (b *Bucket) search(value int) {
	for i := 0; i < b.Size; i++ {
		if b.ValueArr[i] == value {
			log.Printf("Found %v at index %v", value, i)
		}
	}
}

// isEmpty is a Method for indicating if a bucket is empty
func (b *Bucket) isEmpty() bool {
	for i := 0; i < b.Size; i++ {
		if b.ValueArr[i] != 0 {
			return false
		}
	}
	return true
}

// IsFull is a Method for indicating if a bucket is full
func (b *Bucket) isFull() bool {
	for i := 0; i < b.Size; i++ {
		if b.ValueArr[i] == 0 {
			return false
		}
	}
	return true
}

// display is a helper Method for debugging
func (b *Bucket) display() {
	for i := 0; i < b.Size; i++ {
		log.Printf("Index %v: %v", i, b.ValueArr[i])
	}
}

// sort is a helper Method for debugging
func (b *Bucket) sort() {
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size-i-1; j++ {
			if b.ValueArr[j] > b.ValueArr[j+1] {
				b.ValueArr[j], b.ValueArr[j+1] = b.ValueArr[j+1], b.ValueArr[j]
			}

		}
	}
}

// getFirst is a method for getting the first value of a bucket
func (b *Bucket) getFirst() int {
	return b.ValueArr[0]
}

// check is a method for seeing if a value exists ina  bucket
func (b *Bucket) check(value int) bool {
	for i := 0; i < b.Size; i++ {
		if b.ValueArr[i] != value {
			return false
		}
	}

	return true
}
