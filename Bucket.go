package main

import "log"

type Bucket struct {
	ValueArr   []string
	LocalDepth int
	Size       int
}

func NewBucket2(bucketSize int) *Bucket {
	ld := 1
	sz := bucketSize
	valueArr := make([]string, sz)
	for i := 0; i < sz; i++ {
		valueArr[i] = ""
	}
	return &Bucket{
		ValueArr:   valueArr,
		LocalDepth: ld,
		Size:       sz,
	}
}

func NewBucket() *Bucket {
	ld := 1
	sz := 4
	valueArr := make([]string, sz)
	for i := 0; i < sz; i++ {
		valueArr[i] = ""
	}
	return &Bucket{
		ValueArr:   valueArr,
		LocalDepth: ld,
		Size:       sz,
	}
}

func (b *Bucket) getLD() int {
	return b.LocalDepth
}

func (b *Bucket) getSize() int {
	return b.Size
}

func (b *Bucket) getArr() []string {
	return b.ValueArr
}

func (b *Bucket) insert(value string) {
	for i := 0; i < b.Size; i++ {
		if b.ValueArr[i] == "" {
			b.ValueArr[i] = value
			break
		}
	}
}

func (b *Bucket) remove(value string) {
	for i := 0; i < b.Size; i++ {
		if b.ValueArr[i] == value {
			b.ValueArr[i] = ""
			break
		}
	}
}

// find is a function for indicating presence of a value in a bucket
func (b *Bucket) find(value string) bool {
	for i := 0; i < b.Size; i++ {
		if b.ValueArr[i] == value {
			return true
		}
	}
	return false
}

// search is a helper function for debugging
func (b *Bucket) search(value string) {
	for i := 0; i < b.Size; i++ {
		if b.ValueArr[i] == value {
			log.Printf("Found %v at index %v", value, i)
		}
	}
}

func (b *Bucket) isEmpty() bool {
	for i := 0; i < b.Size; i++ {
		if b.ValueArr[i] != "" {
			return false
		}
	}
	return true
}

func (b *Bucket) isFull() bool {
	for i := 0; i < b.Size; i++ {
		if b.ValueArr[i] == "" {
			return false
		}
	}
	return true
}

func (b *Bucket) display() {
	for i := 0; i < b.Size; i++ {
		log.Printf("Index %v: %v", i, b.ValueArr[i])
	}
}

func (b *Bucket) sort() {
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size-i-1; j++ {
			if b.ValueArr[j] > b.ValueArr[j+1] {
				b.ValueArr[j], b.ValueArr[j+1] = b.ValueArr[j+1], b.ValueArr[j]
			}
		}
	}
}

func (b *Bucket) getFirst() string {
	return b.ValueArr[0]
}

func (b *Bucket) check(value string) bool {
	for i := 0; i < b.Size; i++ {
		if b.ValueArr[i] != value {
			return false
		}

	}

	return true
}
