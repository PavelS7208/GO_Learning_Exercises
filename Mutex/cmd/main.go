package main

import "github.com/pavel/msync"

func main() {
	mutex := msync.NewRWPMutex()
	_ = mutex
}
