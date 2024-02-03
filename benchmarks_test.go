package main

import (
	"log"
	"os"
	"testing"
)

func BenchmarkReadFile(b *testing.B) {
	f, err := os.Open("./misc/large.txt")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			//...
		}
	}()

	fileInfo, _ := f.Stat()
	fileSize := fileInfo.Size()

	buf := make([]byte, fileSize)
	log.Print("Iterations count: ", b.N)
	for k := 0; k < b.N; k++ {
		f.ReadAt(buf, 0)
	}
}
