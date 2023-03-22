package main

import (
	"fmt"
	"github.com/faizanaryan94/toolkit"
)

func main() {
	var tools toolkit.Tools
	s := tools.RandomString(10)
	fmt.Println("Random String: ", s)
}
