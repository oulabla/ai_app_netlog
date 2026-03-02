package server

import (
	"fmt"
	"sync"
)

type DependencyInjector struct {
	mu sync.RWMutex
	m  map[string]any
}

var (
	di   *DependencyInjector
	once sync.Once
)

// Singleton с ленивой инициализацией
func GetInjector() *DependencyInjector {
	once.Do(func() {
		di = &DependencyInjector{
			m: make(map[string]any),
		}
	})
	return di
}

// Set - безопасная установка с блокировкой
func (d *DependencyInjector) Set(name string, injection any) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.m[name] = injection
}

// Get - безопасное получение
func (d *DependencyInjector) Get(name string) any {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.m[name]
}

// GetWithType - получение с приведением типа (дженерики)
func GetWithType[T any](name string) (T, error) {
	injector := GetInjector()
	val := injector.Get(name)

	var zero T
	if val == nil {
		return zero, fmt.Errorf("dependency '%s' not found", name)
	}

	typedVal, ok := val.(T)
	if !ok {
		return zero, fmt.Errorf("dependency '%s' has wrong type: expected %T, got %T",
			name, zero, val)
	}

	return typedVal, nil
}

// MustGetWithType - паникует если не найдено
func MustGetWithType[T any](name string) T {
	val, err := GetWithType[T](name)
	if err != nil {
		panic(err)
	}
	return val
}

// Has - проверка существования
func (d *DependencyInjector) Has(name string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, ok := d.m[name]
	return ok
}

// Remove - удаление зависимости
func (d *DependencyInjector) Remove(name string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.m, name)
}

// Clear - очистка всех зависимостей
func (d *DependencyInjector) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.m = make(map[string]any)
}

// Keys - получение всех ключей
func (d *DependencyInjector) Keys() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	keys := make([]string, 0, len(d.m))
	for k := range d.m {
		keys = append(keys, k)
	}
	return keys
}
