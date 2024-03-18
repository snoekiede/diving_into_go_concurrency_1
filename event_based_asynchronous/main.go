package main

import "fmt"

type ResizeEvent struct {
	width  int
	height int
}

type ResizeEventHandler func(event ResizeEvent) (int, int, error)

type ResizeEventListener struct {
	events  chan ResizeEvent
	handler ResizeEventHandler
}

func CreateResizeEventListener(handler ResizeEventHandler) ResizeEventListener {
	return ResizeEventListener{make(chan ResizeEvent), handler}
}

func (l *ResizeEventListener) Start(w *Window) {
	go func() {
		for event := range l.events {
			new_width, new_height, error := l.handler(event)
			if error != nil {
				break
			}
			w.width = new_width
			w.height = new_height
		}
	}()
}

func (l *ResizeEventListener) Stop() {
	close(l.events)
}

func (l *ResizeEventListener) Send(event ResizeEvent) {
	l.events <- event
}

type Window struct {
	listener ResizeEventListener
	title    string
	width    int
	height   int
}

func CreateWindow(title string, width int, height int, handler ResizeEventHandler) Window {
	return Window{CreateResizeEventListener(handler), title, width, height}
}

func (w *Window) Open() {
	w.listener.Start(w)
	fmt.Printf("Window %s opened with size %dx%d\n", w.title, w.width, w.height)
}

func (w *Window) Close() {
	w.listener.Stop()
	fmt.Printf("Window %s closed\n", w.title)
}

func (w *Window) Resize(event ResizeEvent) {
	w.listener.Send(event)
}

func main() {
	window := CreateWindow("My Window", 800, 600, func(event ResizeEvent) (int, int, error) {
		fmt.Printf("Window resized to %dx%d\n", event.width, event.height)
		return event.width, event.height, nil
	})
	window.Open()
	window.Resize(ResizeEvent{1024, 768})
	window.Resize(ResizeEvent{644, 484})
	window.Resize(ResizeEvent{800, 600})
	window.Close()
	fmt.Printf("Height is %d\n", window.height)
	fmt.Printf("Width is %d\n", window.width)

}
