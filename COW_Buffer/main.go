package main

import "fmt"

func main() {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()
	fmt.Println("COWBuffer создан")

	copy1 := buffer.Clone()
	defer copy1.Close()
	copy2 := buffer.Clone()
	copy2.Close()
	copy2.Close()
	copy2.Close()
	_ = copy2.Clone()
	defer copy2.Close()

	fmt.Println(buffer.Update(0, 'g'))  // true
	fmt.Println(buffer.Update(-1, 'g')) // False
	fmt.Println(buffer.Update(4, 'g'))  // False

	data[1] = 'p'
	fmt.Println(buffer.String())
	fmt.Println(copy1.String())
	fmt.Println(copy2.String())

}
