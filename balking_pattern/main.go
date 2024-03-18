package main

import (
	"errors"
	"fmt"
)

type OpenState struct{}
type ClosedState struct{}
type UndefinedState struct{}

type MemoryFile struct {
	name  string
	state interface{}
}

func NewMemoryFile(name string) *MemoryFile {
	return &MemoryFile{name: name, state: UndefinedState{}}
}

func (f *MemoryFile) Open() (*MemoryFile, error) {
	switch f.state.(type) {
	case ClosedState:
		return &MemoryFile{name: f.name, state: OpenState{}}, nil
	case UndefinedState:
		return &MemoryFile{name: f.name, state: OpenState{}}, nil
	default:
		return &MemoryFile{name: f.name, state: UndefinedState{}}, errors.New("file already open")
	}
}

func (f *MemoryFile) Close() (*MemoryFile, error) {
	switch f.state.(type) {
	case OpenState:
		return &MemoryFile{name: f.name, state: ClosedState{}}, nil
	default:
		return &MemoryFile{name: f.name, state: UndefinedState{}}, errors.New("file already closed")
	}
}

func (f *MemoryFile) Read() (*string, error) {
	switch f.state.(type) {
	case OpenState:
		return &f.name, nil
	default:
		return nil, errors.New("file not ready to read from since it's either closed or undefined")
	}
}
func main() {
	file := NewMemoryFile("test.txt")

	file, err := file.Open()
	if err != nil {
		fmt.Printf("Something went wrong opening the file: %v\n", err)
	} else {
		println("file opened")
	}

	content, err := file.Read()
	if err != nil {
		fmt.Printf("Something went wrong reading the file: %v\n", err)
	} else {
		println("file content: ", *content)
	}
	file, err = file.Close()
	if err != nil {
		fmt.Printf("Something went wrong closing the file: %v\n", err)
	} else {
		println("file closed")
	}

	content, err = file.Read()
	if err != nil {
		fmt.Printf("Something went wrong reading the file after closing: %v\n", err)
	} else {
		println("file content: ", *content)
	}

	file, err = file.Close()
	if err != nil {
		fmt.Printf("Something went wrong closing the file: %v\n", err)
	} else {
		println("file closed")
	}
}
