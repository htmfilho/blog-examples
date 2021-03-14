package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

type Publisher interface {
	register(subscriber *Subscriber)
	unregister(subscriber *Subscriber)
	notify(path, event string)
}

type Subscriber interface {
	Receive(path, event string)
}

// PathWatcher observes changes in the file system and works as a Publisher for
// the application by notifying subscribers, which will perform other operations.
type PathWatcher struct {
	subscribers []*Subscriber
	watcher     fsnotify.Watcher
}

func (pw *PathWatcher) observePath(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error", err)
	}
	defer watcher.Close()

	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsDir() {
			return watcher.Add(path)
		}

		return nil
	}); err != nil {
		fmt.Println("ERROR", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				pw.notify(event.Name, event.Op.String())
			case err := <-watcher.Errors:
				fmt.Println("Error", err)
			}
		}
	}()

	<-done
}

func (pw *PathWatcher) register(subscriber *Subscriber) {
	pw.subscribers = append(pw.subscribers, subscriber)
}

func (pw *PathWatcher) unregister(subscriber *Subscriber) {
	length := len(pw.subscribers)

	for i, sub := range pw.subscribers {
		if sub == subscriber {
			pw.subscribers[i] = pw.subscribers[length-1]
			pw.subscribers = pw.subscribers[:length-1]
			break
		}
	}
}

func (pw *PathWatcher) notify(path, event string) {
	for _, sub := range pw.subscribers {
		(*sub).Receive(path, event)
	}
}

func main() {
	pathWatcher := PathWatcher{}

	var pathIndexer Subscriber = &PathIndexer{}
	pathWatcher.register(&pathIndexer)

	var pathFileMD5 Subscriber = &PathFileMD5{}
	pathWatcher.register(&pathFileMD5)

	pathWatcher.observePath("/home/htmfilho/liftbox")
}
