package main

import "fmt"

var (
)

func main() {
	m := make(map[int]int)
	m[1]=1
	m[2]=1
	fmt.Println(len(m))
}
