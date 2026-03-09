package main

import (
	"fmt"
	"unsafe"
)

func main() {

	memory := []byte{
		5, 1, 2, 3, 4, 5, // блок 0 длиной 5
		0, 0, 0, 0, 0, // дырка
		4, 6, 7, 8, 9, // блок 1 длиной 4
		0, 0, 0, 0, 0, 0, // дырка
		6, 11, 12, 13, 14, 15, 16, // блок 2 длиной 6
	}

	pointers := []unsafe.Pointer{
		ptrAt(memory, 1),  // блок 0
		ptrAt(memory, 23), // блок 2
		ptrAt(memory, 12), // блок 1
	}

	fmt.Println(memory)
	Defragment(memory, pointers)
	fmt.Println(memory)

}

// ptrAt - создать указатель на позицию в слайсе (для удобства)
func ptrAt(mem []byte, idx int) unsafe.Pointer {
	return unsafe.Pointer(&mem[idx])
}
