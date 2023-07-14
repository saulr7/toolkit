package main

import (
	"fmt"
	"toolkit/toolkit"
)

func main() {

	var tools toolkit.Tools
	s := tools.RandomString(20)
	fmt.Println(s)

}
