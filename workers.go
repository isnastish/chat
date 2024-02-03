package main

import (
	"io"
	"log"
	"os"

	"github.com/niemeyer/golang/src/pkg/container/vector"
)

/*
Ideas:

Workers participate in file transferring process as peers.
Each worker reads a chunk from a file 256/1024/4096 KB, and uploads them to the client.
Once we've reached the end of file, a confirmation message should be sent to the client that file has been
successfully transmitted.
*/

type Worker struct {
	chunksRead uint32
	chunks     map[int64][]byte
	offset     chan int64
	chunkSize  chan int64
}

type WorkerPool struct {
	workers vector.Vector
}

type Orchestrator struct {
}

func NewWorkerPool(workersCount uint32) *WorkerPool {
	wp := WorkerPool{}
	for i := 0; i < int(workersCount); i++ {
		wp.workers.Push(NewWorker())
	}
	return &wp
}

func NewWorker() *Worker {
	return &Worker{}
}

func (wp *WorkerPool) RegisterWorker(w *Worker) {

}

// NOTE(alx): Pass an offset through channel?
func (w *Worker) ReadChunk(f *os.File, chunkSize int64) {
	for offset := range w.offset {
		buf := make([]byte, 256)
		bytesRead, err := f.ReadAt(buf, offset)
		if err != nil {
			log.Fatal(err)
		} else if err == io.EOF {
			w.chunks[offset] = []byte{}
		}
		log.Println("Bytes read: ", bytesRead)
		w.chunks[offset] = buf
	}
}

/*
Let's assume that we spawn 8 workers, and we have a file, which occupies 4096 bytes on the drive.

Each of them would read 256 (byte) chunk until reaches and of the file. (EOF).

We would have to synchronize work between each worker in order to prevent reading the same chunks.

-- here file starts.
3245989sfds  256 bytes (worker 0) range: [0b - 256b)
3245989sfds  256 bytes (worker 1) range: [256b - 512b)
3245989sfds  256 bytes (worker 3) range: [512b - 768b)
...
3245989sfds (4096)
*/
