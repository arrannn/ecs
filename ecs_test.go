package ecs_test

import (
	"testing"

	"github.com/arrannn/ecs"
)

type Position struct {
	X, Y float64
}

type Velocity struct {
	X, Y float64
}

type Health struct {
	Current int
}

func TestComponentOperations(t *testing.T) {
	e := ecs.NewEngine()

	// Test basic component addition and retrieval
	id1 := e.NewID()
	ecs.Write(e, id1, Position{X: 1, Y: 2})
	ecs.Write(e, id1, Velocity{X: 3, Y: 4})
	ecs.Write(e, id1, Health{Current: 100})

	if pos, exists := ecs.Read[Position](e, id1); !exists || pos.X != 1 || pos.Y != 2 {
		t.Error("Position component not correctly stored/retrieved")
	}
	if vel, exists := ecs.Read[Velocity](e, id1); !exists || vel.X != 3 || vel.Y != 4 {
		t.Error("Velocity component not correctly stored/retrieved")
	}

	// Test component updates
	ecs.Write(e, id1, Position{X: 5, Y: 6})
	if pos, exists := ecs.Read[Position](e, id1); !exists || pos.X != 5 || pos.Y != 6 {
		t.Error("Position component not correctly updated")
	}

	// Test multiple entities
	id2 := e.NewID()
	ecs.Write(e, id2, Position{X: 10, Y: 20})
	ecs.Write(e, id2, Health{Current: 50})

	// Test component deletion
	ecs.Delete[Position](e, id1)
	if _, exists := ecs.Read[Position](e, id1); exists {
		t.Error("Position component should have been deleted")
	}
	if vel, exists := ecs.Read[Velocity](e, id1); !exists || vel.X != 3 || vel.Y != 4 {
		t.Error("Velocity component should still exist and be unchanged")
	}

	// Test DeleteAll
	ecs.DeleteAll(e, id1)
	if _, exists := ecs.Read[Velocity](e, id1); exists {
		t.Error("Velocity component should have been deleted by DeleteAll")
	}
	if _, exists := ecs.Read[Health](e, id1); exists {
		t.Error("Health component should have been deleted by DeleteAll")
	}

	// Verify id2's components still exist
	if pos, exists := ecs.Read[Position](e, id2); !exists || pos.X != 10 || pos.Y != 20 {
		t.Error("Position component for id2 should still exist and be unchanged")
	}

	// Test ForEach
	count := 0
	totalHealth := 0
	ecs.ForEach[Health](e, func(id uint32, health Health) {
		count++
		totalHealth += health.Current
	})
	if count != 1 || totalHealth != 50 {
		t.Error("ForEach didn't iterate correctly over Health components")
	}

	// Test component ordering after deletions
	id3 := e.NewID()
	id4 := e.NewID()
	ecs.Write(e, id3, Position{X: 30, Y: 40})
	ecs.Write(e, id4, Position{X: 50, Y: 60})

	// Delete middle entity
	ecs.Delete[Position](e, id3)

	// Verify id4's position is still correct
	if pos, exists := ecs.Read[Position](e, id4); !exists || pos.X != 50 || pos.Y != 60 {
		t.Error("Position component for id4 should be unchanged after deleting id3")
	}

	// Test non-existent components
	if _, exists := ecs.Read[Velocity](e, id4); exists {
		t.Error("Should not find non-existent component")
	}
}

func TestForEachOrdering(t *testing.T) {
	e := ecs.NewEngine()

	// Add components in specific order
	ids := make([]uint32, 3)
	positions := []Position{
		{X: 1, Y: 1},
		{X: 2, Y: 2},
		{X: 3, Y: 3},
	}

	for i := range positions {
		ids[i] = e.NewID()
		ecs.Write(e, ids[i], positions[i])
	}

	// Delete middle component
	ecs.Delete[Position](e, ids[1])

	// Verify iteration order
	expectedPositions := []Position{
		{X: 1, Y: 1},
		{X: 3, Y: 3}, // Last element should have moved to fill the gap
	}

	index := 0
	ecs.ForEach(e, func(id uint32, pos Position) {
		if pos != expectedPositions[index] {
			t.Errorf("Position at index %d incorrect. Expected %v, got %v",
				index, expectedPositions[index], pos)
		}
		index++
	})

	if index != len(expectedPositions) {
		t.Errorf("Expected %d iterations, got %d", len(expectedPositions), index)
	}
}
