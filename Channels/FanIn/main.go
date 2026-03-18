package main

import (
	"fmt"
	"os"
)

func sum(slice []int, start, end int, res chan int) {
	var sum int
	for _, v := range slice[start:end] {
		sum += v
	}
	res <- sum
}

func run() error {
	array := make([]int, 100)
	for i := 0; i < len(array); i++ {
		array[i] = i + 1
	}

	const numParts = 4
	// Проверка на некорректное количество частей
	if numParts <= 0 {
		return fmt.Errorf("количество частей должно быть больше 0")
	}
	if numParts > len(array) {
		return fmt.Errorf("количество частей (%d) больше длины массива (%d)\n", numParts, len(array))
	}

	// Ожидаем столько результатов, на сколько частей бьем массив
	result := make(chan int, numParts)

	size := len(array) / numParts
	for i := 0; i < numParts; i++ {
		start := i * size
		end := start + size
		if i == numParts-1 {
			end = len(array)
		}
		go sum(array, start, end, result)
	}

	var total int
	for i := 0; i < numParts; i++ {
		sum := <-result
		fmt.Println(sum)
		total += sum
	}

	fmt.Println("Итого:", total)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
