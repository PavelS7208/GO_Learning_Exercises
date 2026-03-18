package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func worker(ctx context.Context, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker останавливаем")
			return
		case v, ok := <-jobs:
			if !ok {
				return
			}
			// Имитация длительной работы (не через Sleep)
			select {
			case <-time.After(2000 * time.Millisecond):
				// работа завершена
			case <-ctx.Done():
				fmt.Println("Worker прерван во время работы")
				return
			}

			// Отправка результата
			select {
			case results <- v * v:
			case <-ctx.Done():
				fmt.Println("Worker прерван при отправке")
				return
			}
		}
	}
}

func input(ctx context.Context, jobs chan<- int) {
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return
		case jobs <- i + 1:
		}
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка сигналов ОС (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nПрерываемся ...")
		cancel()
	}()

	jobs := make(chan int, 10)
	results := make(chan int)

	var wg sync.WaitGroup
	wg.Add(3)

	// Запускаем воркеры и ждем пока данные появятся
	for i := 0; i < 3; i++ {
		go worker(ctx, jobs, results, &wg)
	}

	// Запускаем продюсера
	go func() {
		input(ctx, jobs)
		close(jobs) // Кто создает канал, тот и закрывает
	}()

	// Ждем завершения всех воркеров
	go func() {
		wg.Wait()
		close(results) // Кто создает канал, тот и закрывает
	}()

	// Результат печатаем из канала
	for v := range results {
		fmt.Println(v)
	}
}
