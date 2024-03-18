package main

import (
	"fmt"
	"sync"
)

type ExpensiveCar struct {
	Name  string
	Price uint32
}

type Lazy struct {
	once  sync.Once
	value *ExpensiveCar
	mu    sync.Mutex
}

func (l *Lazy) GetOrInit(initFunc func() *ExpensiveCar) *ExpensiveCar {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.value != nil {
		return l.value
	}

	l.once.Do(func() {
		l.value = initFunc()
	})
	return l.value
}

func main() {
	lazy := Lazy{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		lazy.GetOrInit(func() *ExpensiveCar {
			return &ExpensiveCar{
				Name:  "Expensive Brand",
				Price: 100000,
			}
		})
		fmt.Printf("Value is %+v\n", lazy.value)
	}()

	wg.Wait()

}
