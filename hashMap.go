package main

const InitMapSize = 8

type Node struct {
	key   string
	value int
	next  *Node
}

type HashMap struct {
	size    int
	cap     int
	nodeArr []*Node
}

func NewHashMap() *HashMap {
	return &HashMap{
		size:    0,
		cap:     InitMapSize,
		nodeArr: make([]*Node, InitMapSize),
	}
}

func (h *HashMap) hashFunction(key string) int {
	return len(key) % h.cap
}

func (h *HashMap) Add(key string, value int) {
	if IncrementIfExists(h, key) {
		return
	}
	h.nodeArr[h.hashFunction(key)] = &Node{
		key:   key,
		value: value,
		next:  h.nodeArr[h.hashFunction(key)],
	}
	h.size++
	if h.size >= h.cap/2 {
		h.resize()
	}
}

func IncrementIfExists(h *HashMap, key string) bool {
	value, exists := h.Get(key)
	if exists {
		h.Add(key, value+1)
		return true
	}
	return false
}

func (h HashMap) Get(key string) (int, bool) {
	index := h.hashFunction(key)
	node := h.nodeArr[index]
	for node != nil {
		if node.key == key {
			return node.value, true
		}
		node = node.next
	}
	return 0, false
}

// Delete removes the key-value pair with the specified key from the hash map
func (h *HashMap) Delete(key string) {
	index := h.hashFunction(key)
	node := h.nodeArr[index]
	var prev *Node
	for node != nil {
		if node.key == key {
			if prev == nil {
				h.nodeArr[index] = node.next
			} else {
				prev.next = node.next
			}
			h.size--
			return
		}
		prev = node
		node = node.next
	}
}

// resize resizes the hash map when it becomes too full
func (h *HashMap) resize() {
	h.cap *= 2
	tempNodes := make([]*Node, h.cap)
	for i := 0; i < len(h.nodeArr); i++ {
		node := h.nodeArr[i]
		for node != nil {
			index := h.hashFunction(node.key)
			newNode := &Node{
				key:   node.key,
				value: node.value,
				next:  tempNodes[index],
			}
			tempNodes[index] = newNode
			node = node.next
		}
	}
	h.nodeArr = tempNodes
}
