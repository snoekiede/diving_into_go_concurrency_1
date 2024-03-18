package main

import (
	"fmt"
	"sync"
)

type Stock struct {
	Name  string
	Price float64
}

type Monitor struct {
	Value     *Stock
	Mutex     sync.Mutex
	Condition *sync.Cond
}

func NewStock(name string, price float64) *Stock {
	return &Stock{
		Name:  name,
		Price: price,
	}
}

func NewMonitor(stock *Stock) *Monitor {
	return &Monitor{
		Value:     stock,
		Mutex:     sync.Mutex{},
		Condition: sync.NewCond(&sync.Mutex{}),
	}
}

func (m *Monitor) UpdatePrice(price float64) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	fmt.Printf("Updating price of %v to %v for stock %v\n", m.Value.Name, price, m.Value.Name)
	m.Value.Price = price
	m.Condition.Broadcast()
}

func (m *Monitor) WaitForRelease() {
	limit := 120.0
	m.Mutex.Lock()
	for m.Value.Price < limit {
		fmt.Printf("Waiting for price of %v to drop below %v\n", m.Value.Name, limit)
		m.Condition.Wait()
	}
	fmt.Printf("Price of %v is now %v and above %v\n", m.Value.Name, m.Value.Price, limit)
	m.Mutex.Unlock()
}

func main() {
	monitor := NewMonitor(NewStock("GOOG", 110.0))
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			monitor.UpdatePrice(110.0 + float64(i)*2.0)
		}()
	}

	go monitor.WaitForRelease()
	wg.Wait()
	fmt.Printf("Stock is now %+v\n", monitor.Value)
}
