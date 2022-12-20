package inter

import "fmt"

type Human interface {
	printName()
}

type Person struct {
	name    string
	address int64
}

func (p *Person) init(name string, address int64) {
	p.name = name
	p.address = address
}

func (p *Person) printName() {
	fmt.Println("person")
}

func PrintMsg(human Human) {
	human.printName()
}
