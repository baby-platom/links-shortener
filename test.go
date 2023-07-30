package main

import "fmt"

type A struct {
	Data map[int]string
}

func (a *A) Add(key int, value string) {
	a.Data[key] = value
}

type B struct {
	A
}

func (b *B) Add(key int, value string) {
	key += 10
	b.A.Add(key, value)
}

func main() {
	b := &B{A{Data: make(map[int]string)}}
	b.Add(5, "hello")
	fmt.Println(b.Data)
}
