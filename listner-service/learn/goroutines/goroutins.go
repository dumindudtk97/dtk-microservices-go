package main

import (
	"fmt"
	"math/rand"
	"time"
)

//  goroutine is a function that is capable of running concurrently with other functions

// function to pring 0-9 with sleep inbetween
func f(n int) {
	for i := 0; i < 10; i++ {
		fmt.Println(n, ":", i)
		amt := time.Duration(rand.Intn(250))
		time.Sleep(time.Millisecond * amt)
	}
}

func main() {
	// call go routines // they all run simultaneously
	for i := 0; i < 10; i++ {
		go f(i)
	}
	var input string
	fmt.Scanln(&input) // to prevent main from exiting before go routines finish
}
