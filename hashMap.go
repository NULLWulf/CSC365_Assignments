package main

import (
	"math/rand"
)

const initialSize = 8 // initial size of the hash map

// Node represents a key-value pair stored in the hash map
type Node struct {
	key    string
	values []string
	next   *Node
}

// HashMap represents the hash map data structure
type HashMap struct {
	size     int
	capacity int
	buckets  []*Node
}

// NewHashMap creates a new hash map with the specified initial size
func NewHashMap() *HashMap {
	return &HashMap{
		size:     0,
		capacity: initialSize,
		buckets:  make([]*Node, initialSize),
	}
}

// hashFunction maps the key to an index in the hash map
func (hm *HashMap) hashFunction(key string) int {
	return len(key) % hm.capacity
}

// Add adds a new string value to the array of values associated with the specified key in the hash map
func (hm *HashMap) Add(key string, value string) {
	index := hm.hashFunction(key)
	node := hm.buckets[index]
	for node != nil {
		if node.key == key {
			node.values = append(node.values, value)
			return
		}
		node = node.next
	}
	newNode := &Node{
		key:    key,
		values: []string{value},
		next:   hm.buckets[index],
	}
	hm.buckets[index] = newNode
	hm.size++
	if hm.size >= hm.capacity/2 {
		hm.resize()
	}
}

// Get retrieves a random string value associated with the specified key from the hash map
func (hm *HashMap) Get(key string) (string, bool) {
	index := hm.hashFunction(key)
	node := hm.buckets[index]
	for node != nil {
		if node.key == key {
			if len(node.values) == 0 {
				return "", false
			}
			randomIndex := rand.Intn(len(node.values))
			return node.values[randomIndex], true
		}
		node = node.next
	}
	return "", false
}

// Delete removes the key-value pair with the specified key from the hash map
func (hm *HashMap) Delete(key string) {
	index := hm.hashFunction(key)
	node := hm.buckets[index]
	var prev *Node
	for node != nil {
		if node.key == key {
			if prev == nil {
				hm.buckets[index] = node.next
			} else {
				prev.next = node.next
			}
			hm.size--
			return
		}
		prev = node
		node = node.next
	}
}

// resize resizes the hash map when it becomes too full
func (hm *HashMap) resize() {
	hm.capacity *= 2
	newBuckets := make([]*Node, hm.capacity)
	for i := 0; i < len(hm.buckets); i++ {
		node := hm.buckets[i]
		for node != nil {
			index := hm.hashFunction(node.key)
			newNode := &Node{
				key:    node.key,
				values: node.values,
				next:   newBuckets[index],
			}
			newBuckets[index] = newNode
			node = node.next
		}
	}
	hm.buckets = newBuckets
}
