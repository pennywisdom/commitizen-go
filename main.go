package main

import (
	"github.com/pennywisdom/commitizen-go/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
