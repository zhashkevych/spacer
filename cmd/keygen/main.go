package main

import (
	"fmt"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" + "!@#$%^&*()_+}:><?"

func main() {
	rand.Seed(time.Now().Unix())

	key := make([]byte, 32)
	for i := range key {
		key[i] = charset[rand.Intn(len(charset))]
	}

	fmt.Println("New encryption key:", string(key))
}