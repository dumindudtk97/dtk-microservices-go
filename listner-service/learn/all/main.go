package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 2)     // buffer size = 2
	exit := make(chan struct{}) // used to do exit on main (once all go routines completed)

	go func() {
		for i := 0; i < 5; i++ { // send 0,1,2,3,4 to channel ch with 1 second interval

			fmt.Println(time.Now(), i, "sending")
			ch <- i
			fmt.Println(time.Now(), i, "sent")

			time.Sleep(1 * time.Second)
		}

		fmt.Println(time.Now(), "all completed, leaving")

		close(ch)
	}()

	go func() {
		// XXX: This is overcomplicated because is only channel only, "select"
		// shines when using multiple channels.
		for {
			select {
			case v, open := <-ch: // getting channel status and value
				if !open {
					close(exit) // if ch is closed close exit // ch is closed after sending last number, so it and all completed, leaving is already printed
					return
				}

				fmt.Println(time.Now(), "received", v)
			}
		}

		// XXX: In cases where only one channel is used
		// for v := range ch {
		// 	fmt.Println(time.Now(), "received", v)
		// }

		// close(exit)
	}()

	fmt.Println(time.Now(), "waiting for everything to complete") // this line prints first

	<-exit // when exit closed this not exist, main can move on

	fmt.Println(time.Now(), "exiting")
}
