package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

const (
	charset  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+}:><?"
	filename = "enc.key"
)

func main() {
	rand.Seed(time.Now().Unix())

	key := make([]byte, 32)
	for i := range key {
		key[i] = charset[rand.Intn(len(charset))]
	}

	if err := ioutil.WriteFile(filename, key, 0777); err != nil {
		log.Fatal(err)
	}

	fmt.Println("New encryption key generated in", filename)
}
