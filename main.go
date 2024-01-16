package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func downloadWebpage(workerGroup *sync.WaitGroup, id int, url string) {
	fmt.Printf("Worker %d started\n", id)
	response, _ := http.Get(url)
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	filename := fmt.Sprintf("file_%d.html", id)
	os.WriteFile(filename, body, 0644)
	workerGroup.Done() // Decrement the WaitGroup counter when the goroutine completes
	fmt.Printf("Worker %d completed\n", id)
}

func main() {
	var workerGroup sync.WaitGroup

	urls := []string{
		"https://example.com",
		"https://www.google.com/",
		"https://store.steampowered.com/",
		"https://pomodorotimer.online/",
	}

	for i, url := range urls {
		// Increment the WaitGroup counter for each goroutine
		workerGroup.Add(1)
		go downloadWebpage(&workerGroup, i, url)
	}

	workerGroup.Wait()

	fmt.Println("Finished")
}
