package main

import "fmt"

func generator(n int) <-chan int {
	out := make(chan int)

	go func() {
		for i := 0; i < n; i++ {
			out <- i + 1
		}
		close(out)
	}()

	return out
}

func main() {

	ch := generator(13)
	for num := range ch {
		fmt.Println(num)
	}
}
