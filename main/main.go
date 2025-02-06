package main

import (
	"fmt"
	"github.com/dastoikov/cool-lex-go/coollex"
	"log"
)

func main() {
	alg, err := coollex.NewLinkedList(3, 2)
	if err != nil {
		log.Fatal(err)
	}
	for combination := range alg.Combinations() {
		for element := range combination {
			fmt.Print(element)
		}
		fmt.Println()
	}
}
