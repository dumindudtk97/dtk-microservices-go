package main

import (
	"fmt"
	"time"
)

// Channels provide a way for two goroutines to communicate with one another and synchronize their execution.
// <- (left arrow) operator is used to send and receive messages on the channe

// ping forever
func pinger(c chan string) {
	for i := 0; ; i++ {
		c <- "ping" // c <- "ping" means send "ping"
	}
}

// specify direction        here pinger can only send to channel
func pingerSendToChannelOnly(c chan<- string) {}

// pong forever
func ponger(c chan string) {
	for i := 0; ; i++ {
		c <- "pong" // c <- "ping" means send "pong"
	}
}

// recieve a message from c every one second and print
func printer(c chan string) {
	for {
		fmt.Println(<-c) // recieve a message from c   //msg := <-c		 // msg := <- c means receive a message and store it in msg
		time.Sleep(time.Second * 1)
	}
}

// can only recieve from channel
func printerRecieveOnly(c <-chan string) {}

// Using a channel synchronizes the goroutines.
// When pinger attempts to send a message on the channel it will wait until printer is ready to receive the message.
func main() {
	var c chan string = make(chan string)

	go pinger(c)
	go ponger(c)
	go printer(c)

	var input string
	// main exits after this line (go routines terminate also)
	fmt.Scanln(&input)
}

// here pinger and ponger will take turns sending messages
