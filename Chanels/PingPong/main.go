package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func ping(ball chan string, exit <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	ping := "Ping"
	for {
		select {
		case <-exit:
			fmt.Println("Выход из Ping")
			return
		case ball <- ping:
			select {
			case <-exit:
				fmt.Println("Выход из Ping")
				return
			case pong := <-ball:
				fmt.Printf("%s\n", pong)
				time.Sleep(time.Second)
			}
		}
	}
}

func pong(ball chan string, exit <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	pong := "Pong"
	for {
		select {
		case <-exit:
			fmt.Println("Выход из Pong")
			return
		case ping := <-ball:
			fmt.Printf("%s - ", ping)
			time.Sleep(time.Second)
			select {
			case <-exit:
				fmt.Println("Выход из Pong")
				return
			case ball <- pong:
			}
		}
	}
}

// Ждем сигнала от ОС и закрываем канал (пойдет сообщение всем кто слушает канал)
func listenSignal(exit chan<- struct{}) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	fmt.Println("\nПолучен сигнал выхода")
	close(exit)
}

func main() {

	fmt.Println("Начали ping-pong. Ctrl-C выход")

	var wg sync.WaitGroup
	wg.Add(2)

	ball := make(chan string)
	exit := make(chan struct{})
	go ping(ball, exit, &wg)
	go pong(ball, exit, &wg)
	go listenSignal(exit)

	wg.Wait()
	fmt.Println("Завершили ping-pong.")
}
