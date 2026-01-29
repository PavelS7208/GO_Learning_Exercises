package main

import "fmt"

func main() {
	fmt.Println("Начали")
	data := NewOrderedMap[int, int]()

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	fmt.Println(data.Size())
	fmt.Println(data.String())

	//  Проверим посчитать сумму
	var sum int
	data.ForEachInOrder(func(_ int, v int) {
		sum += v
	})
	fmt.Println(sum)
}
