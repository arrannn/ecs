package ecs

import "fmt"

type componentStorage[T any] struct {
	components []T
	idToIndex  map[uint32]int
	indexToID  []uint32
}

func newComponentStorage[T any]() *componentStorage[T] {
	return &componentStorage[T]{components: make([]T, 0), idToIndex: make(map[uint32]int)}
}

func (c *componentStorage[T]) Read(id uint32) (T, bool) {
	if idx, exists := c.idToIndex[id]; exists {
		return c.components[idx], true
	}
	var t T
	return t, false
}

func (c *componentStorage[T]) Write(id uint32, component T) {
	if idx, exists := c.idToIndex[id]; exists {
		c.components[idx] = component
		return
	}
	c.idToIndex[id] = len(c.components)
	c.components = append(c.components, component)
	c.indexToID = append(c.indexToID, id)
}

func (c *componentStorage[T]) Delete(id uint32) {
	idx, exists := c.idToIndex[id]
	if !exists {
		return
	}
	lastIdx := len(c.components) - 1
	if idx < lastIdx {
		lastID := c.indexToID[lastIdx]
		c.components[idx] = c.components[lastIdx]
		c.indexToID[idx] = lastID
		c.idToIndex[lastID] = idx
	}
	c.components = c.components[:lastIdx]
	c.indexToID = c.indexToID[:lastIdx]
	delete(c.idToIndex, id)
}

func (c *componentStorage[T]) String() string {
	result := "" //fmt.Sprintf("  Components (%d):\n", len(c.components))
	for i, comp := range c.components {
		entityID := c.indexToID[i]
		result += fmt.Sprintf("    Entity %d: %+v\n", entityID, comp)
	}
	return result
}
