package main

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru"
)

func main() {
	l, _ := lru.New(10)
	for i := 0; i < 100; i++ {
		l.Add(i, i)
	}
	for l.Len() > 0 {
		fmt.Println(l.RemoveOldest())
	}

}
