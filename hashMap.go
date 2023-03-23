package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
)

const initialSize = 8 // initial Size of the hash map

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

// NewHashMap creates a new hash map with the specified initial Size
func NewHashMap() *HashMap {
	return &HashMap{
		size:     0,
		capacity: initialSize,
		buckets:  make([]*Node, initialSize),
	}
}
func (hm *HashMap) hashFunction(key string) (hash int) {
	hash = 0
	for i := 0; i < len(key); i++ {
		hash = (hash*31 + int(key[i])) % hm.capacity
	}
	return hash
}

func (hm *HashMap) getIndex(key string) (index int) {
	return hm.hashFunction(key) % hm.capacity
}

// Add adds a new string value to the array of values associated with the specified key in the hash map
func (hm *HashMap) Add(key string, value string) {
	index := hm.getIndex(key)
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

func (hm *HashMap) GetKeyList() []string {
	keyList := make([]string, 0)
	for i := 0; i < hm.capacity; i++ {
		node := hm.buckets[i]
		for node != nil {
			keyList = append(keyList, node.key)
			node = node.next
		}
	}
	return keyList
}

// PrintKeys prints all keys in the hash map
func (hm *HashMap) PrintKeys() {
	for i := 0; i < hm.capacity; i++ {
		node := hm.buckets[i]
		for node != nil {
			log.Println(node.key)
			node = node.next
		}
	}
}

func (hm *HashMap) PrintValues() {
	for i := 0; i < hm.capacity; i++ {
		node := hm.buckets[i]
		for node != nil {
			log.Println(node.values)
			node = node.next
		}
	}
}

func (hm *HashMap) SaveToFile(filename string) error {
	// Convert the hashmap to a map[string][]string to make it easier to encode to JSON
	data := make(map[string][]string)
	for _, node := range hm.buckets {
		for node != nil {
			data[node.key] = node.values
			node = node.next
		}
	}

	// Encode the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write the data to the file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadHashMapFromFile(filename string) (*HashMap, error) {
	// Read the JSON data from the file
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Decode the JSON data into a map[string][]string
	data := make(map[string][]string)
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	// Create a new HashMap and populate it with the data
	hm := NewHashMap()
	for key, values := range data {
		for _, value := range values {
			hm.Add(key, value)
		}
	}

	return hm, nil
}
