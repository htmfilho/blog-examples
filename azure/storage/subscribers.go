package main

import "fmt"

type PathIndexer struct {
}

func (pi *PathIndexer) Receive(path, event string) {
	fmt.Printf("Indexing: %v, %v\n", path, event)
}

type PathFileMD5 struct {
}

func (pfm *PathFileMD5) Receive(path, event string) {

	fmt.Printf("Syncing: %v, %v\n", path, event)
}
