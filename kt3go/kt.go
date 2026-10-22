
package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sync"
)

const workers = 4

// Задача 1. Конкурентное вычисление квадратов.
func task1() {
	numbers := []int{2, 4, 6, 8, 10}
	var wg sync.WaitGroup

	for _, n := range numbers {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()
			fmt.Println(n * n)
		}(n)
	}

	wg.Wait()
}

// Задача 2. Конкурентное вычисление суммы квадратов.
func task2() {
	numbers := []int{2, 4, 6, 8, 10}
	results := make(chan int)

	for _, n := range numbers {
		go func(n int) {
			results <- n * n
		}(n)
	}

	sum := 0
	for range numbers {
		sum += <-results
	}

	fmt.Println("Сумма квадратов:", sum)
}

// Задача 3. Конкурентная запись в map.
func task3() {
	numbers := []int{2, 4, 6, 8, 10}
	data := make(map[int]int)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, n := range numbers {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			mu.Lock()
			data[n] = n * n
			mu.Unlock()
		}(n)
	}

	wg.Wait()
	fmt.Println(data)
}

// Задача 4. Конкурентный счётчик.
type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func task4() {
	var counter Counter
	var wg sync.WaitGroup

	const goroutines = 10
	const increments = 1000

	for i := 0; i < goroutines; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < increments; j++ {
				counter.Increment()
			}
		}()
	}

	wg.Wait()
	fmt.Println("Итоговое значение:", counter.Value())
}

// Задача 5. Обработка изображений.
func processImage(filename string) (string, bool) {
	info, err := os.Stat(filename)
	if err != nil || !info.Mode().IsRegular() {
		return "", false
	}

	file, err := os.Open(filename)
	if err != nil {
		return "", false
	}

	config, _, err := image.DecodeConfig(file)
	file.Close()

	if err != nil {
		return "", false
	}

	tag := fmt.Sprintf(
		`<img src="%s" width="%d" height="%d" />`,
		filepath.Base(filename),
		config.Width,
		config.Height,
	)

	return tag, true
}

func imageWorker(jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for filename := range jobs {
		if tag, ok := processImage(filename); ok {
			results <- tag
		}
	}
}

func task5(filenames []string) {
	if len(filenames) == 0 {
		return
	}

	jobs := make(chan string)
	results := make(chan string, len(filenames))

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go imageWorker(jobs, results, &wg)
	}

	go func() {
		for _, filename := range filenames {
			jobs <- filename
		}
		close(jobs)

		wg.Wait()
		close(results)
	}()

	for tag := range results {
		fmt.Println(tag)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run main.go <номер задачи> [файлы]")
		return
	}

	switch os.Args[1] {
	case "1":
		task1()
	case "2":
		task2()
	case "3":
		task3()
	case "4":
		task4()
	case "5":
		task5(os.Args[2:])
	default:
		fmt.Println("Укажите номер задачи от 1 до 5")
	}
}
