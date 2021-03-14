package main

import "fmt"

type PathIndexer struct{}

func (pi *PathIndexer) receive(path, event string) {
	fmt.Printf("Indexing: %v, %v\n", path, event)
}

type PathFileMD5 struct{}

func (pfm *PathFileMD5) receive(path, event string) {
	fmt.Printf("Checksuming: %v, %v\n", path, event)
}
