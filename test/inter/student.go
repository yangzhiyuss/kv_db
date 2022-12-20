package inter

import (
	"fmt"
)

type Student struct {
	Person
}

func (s *Student) printName() {
	fmt.Println("student")
}
