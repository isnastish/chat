package main

import (
// "crypto/rand"
// "log"
// "os"
// "time"
)

// func GenerateLargeFile() {
// 	f, err := os.Create("large.txt")
// 	if err != nil {
// 		log.Fatal("Failed to create file.", err)
// 	}

// 	defer func() {
// 		if err := f.Close(); err != nil {
// 			log.Println("Closed file.")
// 		}
// 	}()

// 	rand.Seed(time.Now().UnixNano())

// 	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 	genStr := func(bytes int) {
// 		b := make([]rune, n)
// 		for i := range b {
// 			b[i] = letterRunes[rand.Intn(len(letterRunes))]
// 		}
// 		return string(b)
// 	}

// 	for i := 0; i < 20000; i++ {

// 	}
// }
