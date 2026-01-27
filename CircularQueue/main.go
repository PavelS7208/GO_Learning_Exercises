package main

import "fmt"

func main() {

	queue := NewCircularQueue[int8](4)
	fmt.Println(queue)

	queue.Push(1)
	fmt.Println(queue)

	queue.Push(2)
	fmt.Println(queue)

	queue.Push(3)
	fmt.Println(queue)

	queue.Pop()
	fmt.Println(queue)

	queue.Push(4)
	fmt.Println(queue)
	if !queue.Push(5) {
		fmt.Println("Не шмагла вставить")
	}
	fmt.Println(queue)

	queue.Pop()
	fmt.Println(queue)

	queue.Pop()
	fmt.Println(queue)

	queue.Pop()
	fmt.Println(queue)

	queue.Pop()
	fmt.Println(queue)

	queue.Push(6)
	fmt.Println(queue)

}
