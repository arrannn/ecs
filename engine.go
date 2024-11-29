package ecs

import (
	"fmt"
	"reflect"
)

type Engine struct {
	stores map[string]any
	nextID uint32
}

func NewEngine() *Engine {
	return &Engine{stores: make(map[string]any), nextID: 0}
}

func (e *Engine) NewID() uint32 {
	id := e.nextID
	e.nextID++
	return id
}

func (e *Engine) String() string {
	result := "Engine State:\n"
	for typeName, store := range e.stores {
		result += fmt.Sprintf("%s:\n", typeName)
		result += store.(fmt.Stringer).String()
	}
	result += fmt.Sprintf("Next ID: %d\n", e.nextID)
	return result
}

func Read[T any](e *Engine, id uint32) (T, bool) {
	store := getOrCreateStore[T](e)
	return store.Read(id)
}

func Write[T any](e *Engine, id uint32, component T) {
	store := getOrCreateStore[T](e)
	store.Write(id, component)
}

func Delete[T any](e *Engine, id uint32) {
	store := getOrCreateStore[T](e)
	store.Delete(id)
}

func DeleteAll(e *Engine, id uint32) {
	for _, store := range e.stores {
		deleteFunc := reflect.ValueOf(store).MethodByName("Delete")
		if deleteFunc.IsValid() {
			deleteFunc.Call([]reflect.Value{reflect.ValueOf(id)})
		}
	}
}

func ForEach[T any](e *Engine, f func(uint32, T)) {
	store := getOrCreateStore[T](e)
	for i, component := range store.components {
		f(store.indexToID[i], component)
	}
}

func getOrCreateStore[T any](e *Engine) *componentStorage[T] {
	typeString := reflect.TypeOf((*T)(nil)).Elem().String()
	if store, exists := e.stores[typeString]; exists {
		return store.(*componentStorage[T])
	}
	store := newComponentStorage[T]()
	e.stores[typeString] = store
	return store
}
