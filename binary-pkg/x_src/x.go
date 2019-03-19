package x

import "fmt"

type Person struct {
	Name string
}

func Add(x, y int) int {
	return x + y
}

func (self Person) Hello() {
	fmt.Printf("Hello %s\n", self.Name)
}
