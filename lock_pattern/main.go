package main

import (
	"fmt"
	"sync"
)

type Version struct {
	Version string
	Content string
}

func NewVersion(version, content string) Version {
	return Version{
		Version: version,
		Content: content,
	}
}

func main() {
	list_lock := sync.Mutex{}

	var versions []Version
	var wg sync.WaitGroup

	for counter := 0; counter < 10; counter++ {
		wg.Add(1)
		go func(counter int) {
			defer wg.Done()
			list_lock.Lock()
			defer list_lock.Unlock()
			versions = append(versions, NewVersion(fmt.Sprintf("v0.%d", counter), fmt.Sprintf("content %d", counter)))
		}(counter)
	}

	wg.Wait()
	fmt.Println("Result: ", versions)
}
