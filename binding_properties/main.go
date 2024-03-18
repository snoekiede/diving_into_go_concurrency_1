package main

import (
	"fmt"
	"sync"
)

type Observable interface {
	PropertyChanged(propertyName string, newValue interface{})
}

type Person struct {
	name      string
	age       int
	observers []Observable
	mu        sync.Mutex
}

func NewPerson(name string, age int) *Person {
	return &Person{
		name: name,
		age:  age,
	}
}

func (p *Person) GetName() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.name
}

func (p *Person) SetName(name string) {
	p.mu.Lock()
	p.name = name
	p.mu.Unlock()
	p.notifyObservers("name", name)
}

func (p *Person) GetAge() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.age
}

func (p *Person) SetAge(age int) {
	p.mu.Lock()
	p.age = age
	p.mu.Unlock()
	p.notifyObservers("age", age)
}

func (p *Person) Subscribe(observer Observable) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.observers = append(p.observers, observer)
}

func (p *Person) notifyObservers(property string, value interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, observer := range p.observers {
		observer.PropertyChanged(property, value)
	}
}

type NameObserver struct{}
type AgeObserver struct{}

func (no *NameObserver) PropertyChanged(propertyName string, newValue interface{}) {
	if propertyName == "name" {
		fmt.Printf("NameObserver: %s changed to %s\n", propertyName, newValue)
	}
}

func (ao *AgeObserver) PropertyChanged(propertyName string, newValue interface{}) {
	if propertyName == "age" {
		age, ok := newValue.(int)
		if ok {
			fmt.Printf("AgeObserver: %s changed to %d\n", propertyName, age)
		}
	}
}

func main() {
	person := NewPerson("Test", 55)

	person.Subscribe(&NameObserver{})
	person.Subscribe(&AgeObserver{})

	go func() {
		person.SetName("Jane")
		person.SetAge(21)

	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		person.SetName("John")
		person.SetAge(55)
	}()
	person.SetAge(32)
	wg.Wait()
}
