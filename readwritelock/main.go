package main

import (
	"fmt"
	"sync"
)

// Version represents a version with content
type Version struct {
	version string
	content string
}

// NewVersion creates a new Version with the given version and content
func NewVersion(version, content string) *Version {
	return &Version{
		version: version,
		content: content,
	}
}

func (v *Version) String() string {
	return fmt.Sprintf("Version: %s, Content: %s", v.version, v.content)
}

func main() {
	var wg sync.WaitGroup
	list := make([]*Version, 0)
	mu := sync.RWMutex{}

	for counter := 0; counter < 10; counter++ {
		wg.Add(1)
		go func(counter int) {
			defer wg.Done()
			mu.Lock()
			version := NewVersion(fmt.Sprintf("v0.%d", counter), fmt.Sprintf("content %d", counter*2))
			list = append(list, version)
			mu.Unlock()
		}(counter)
	}

	wg.Wait()

	mu.RLock()
	fmt.Printf("Result: %v\n", list)
	mu.RUnlock()
}
