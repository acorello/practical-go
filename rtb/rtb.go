package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func main() {
	// We have 50 msec to return an answer
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	url := "https://go.dev" // return the 7¢ ad
	// url := "http://go.dev" // return the default ad
	bid := bidOn(ctx, url)
	fmt.Println(bid)
}

// If algo didn't finish in time, return a default bid
func bidOn(ctx context.Context, url string) Bid {
	//	CONTROLLER
	//		- UNBUFFERED CHANNEL*
	//			- 1 ASYNC WRITER*
	//		- (1 READER* |X| 1 TIMEOUT [+async DRAINER])
	//	UNBUFFERED => couples writer and reader
	//	TIMEOUT either READER or TIMEOUT
	//		IF NO READER, THEN WRITER IS STUCK
	bestBidChan := make(chan Bid)
	go func() {
		//	WORKER
		//		has ref to channel
		//		is the only writer
		defer close(bestBidChan)    // being the only writer it should close the channel
		bestBidChan <- bestBid(url) // it may hang if no readers
	}()
	select {
	case bid := <-bestBidChan:
		//	CONTROLLER.OK
		//		we've read the only result => WORKER terminates
		return bid
	case <-ctx.Done():
		//	CONTROLLER.TIMEOUT
		//		too late, but I have to consume from channel or the writer is stuck
		go func() {
			//	DRAINER
			<-bestBidChan
		}()
		return defaultBid
	}
}

var defaultBid = Bid{
	AdURL: "http://adsЯus.com/default",
	Price: 3,
}

// Written by Algo team, time to completion varies
func bestBid(url string) Bid {
	// Simulate work
	d := 100 * time.Millisecond
	if strings.HasPrefix(url, "https://") {
		d = 20 * time.Millisecond
	}
	time.Sleep(d)

	return Bid{
		AdURL: "http://adsЯus.com/ad17",
		Price: 7,
	}
}

type Bid struct {
	AdURL string
	Price int // In ¢
}
