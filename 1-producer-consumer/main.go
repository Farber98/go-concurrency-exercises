//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

func producer(stream Stream, tweetCh chan<- *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(tweetCh) // As we receive channel ownership, must close it.
			return
		}
		tweetCh <- tweet
	}
}

func consumer(tweets <-chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tweets { //  Exits when channel closes
		if t.IsTalkingAboutGo() {
			log.Println(t.Username, "\ttweets about golang")
		} else {
			log.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	var wg sync.WaitGroup

	start := time.Now()
	stream := GetMockStream()

	tweetCh := make(chan *Tweet)

	wg.Add(1)
	go producer(stream, tweetCh, &wg) // One producer
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go consumer(tweetCh, &wg) // N independent consumers
	}
	wg.Wait() // Sync phase

	fmt.Printf("Process took %s\n", time.Since(start))
}
