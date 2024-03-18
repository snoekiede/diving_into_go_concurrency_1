//package main
//
//import (
//	"fmt"
//	"sync"
//	"time"
//)
//
//var wg sync.WaitGroup
//
//type ExpensiveCar struct {
//	brand string
//	color string
//}
//
//func CreateCar(brand string, color string) *ExpensiveCar {
//	return &ExpensiveCar{brand: brand, color: color}
//}
//
//type GuardedGarage struct {
//	cars chan *ExpensiveCar
//}
//
//func NewGuardedGarage() *GuardedGarage {
//	return &GuardedGarage{
//		cars: make(chan *ExpensiveCar),
//	}
//}
//
//func (g *GuardedGarage) AddCar(car *ExpensiveCar) {
//	fmt.Printf("Adding car: %v\n", car)
//	g.cars <- car
//}
//
//func (g *GuardedGarage) GetCar() *ExpensiveCar {
//	defer wg.Done()
//	for {
//		select {
//		case car := <-g.cars:
//			return car
//		default:
//			fmt.Println("Waiting for car")
//		}
//	}
//}
//
//func main() {
//	garage := NewGuardedGarage()
//
//	garage.AddCar(CreateCar("BMW", "red"))
//	garage.AddCar(CreateCar("Audi", "blue"))
//	garage.AddCar(CreateCar("Mercedes", "green"))
//	wg.Add(1)
//	go func() {
//		for {
//			car := garage.GetCar()
//			fmt.Printf("Got car: %v\n", car)
//		}
//	}()
//	wg.Wait()
//
//}

package main

import (
	"fmt"
)

type ExpensiveCar struct {
	brand string
	color string
}

func CreateCar(brand string, color string) *ExpensiveCar {
	return &ExpensiveCar{brand: brand, color: color}
}

type GuardedGarage struct {
	ch chan *ExpensiveCar
}

func NewGuardedGarage() *GuardedGarage {
	g := &GuardedGarage{}
	g.ch = make(chan *ExpensiveCar)
	return g
}

func (g *GuardedGarage) Put(c *ExpensiveCar) {
	g.ch <- c
	fmt.Println("Putting", c)
}

func (g *GuardedGarage) Get() *ExpensiveCar {
	return <-g.ch
}

func main() {
	garage := NewGuardedGarage()
	go func() {
		garage.Put(CreateCar("BMW", "red"))
		garage.Put(CreateCar("Audi", "blue"))
		garage.Put(CreateCar("Mercedes", "green"))
		close(garage.ch)
	}()

	fmt.Println(garage.Get())
	fmt.Println(garage.Get())
	fmt.Println(garage.Get())
}
