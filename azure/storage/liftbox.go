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
	observe()
}

type Subscriber interface {
	receive(path, event string)
}

// PathWatcher observes changes in the file system and works as a Publisher for
// the application by notifying subscribers, which will perform other operations.
type PathWatcher struct {
	subscribers []*Subscriber
	watcher     fsnotify.Watcher
	rootPath    string
}

// register subscribers to the publisher
func (pw *PathWatcher) register(subscriber *Subscriber) {
	pw.subscribers = append(pw.subscribers, subscriber)
}

// unregister subscribers from the publisher
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

// notify subscribers that a event has happened, passing the path and the type
// of event as message.
func (pw *PathWatcher) notify(path, event string) {
	for _, sub := range pw.subscribers {
		(*sub).receive(path, event)
	}
}

// observe changes to the file system using the fsnotify library
func (pw *PathWatcher) observe() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error", err)
	}
	defer watcher.Close()

	if err := filepath.Walk(pw.rootPath, func(path string, info os.FileInfo, err error) error {
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

func main() {
	var pathWatcher Publisher = &PathWatcher{
		rootPath: "/home/username/liftbox",
	}

	var pathIndexer Subscriber = &PathIndexer{}
	pathWatcher.register(&pathIndexer)

	var pathFileMD5 Subscriber = &PathFileMD5{}
	pathWatcher.register(&pathFileMD5)

	pathWatcher.observe()
}
