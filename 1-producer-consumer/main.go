//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, tweets chan *Tweet, s *sync.WaitGroup) {
	defer s.Done()
	for {
		tweet, err := stream.Next()
		if errors.Is(err, ErrEOF) {
			fmt.Println("EOF")
			break
		}
		tweets <- tweet
	}
	fmt.Println("DONE PRODUCING")
	close(tweets)
}

func consumer(tweets chan *Tweet, s *sync.WaitGroup) {
	defer s.Done()
loop:
	for {
		select {
		case t, ok := <-tweets:
			if !ok {
				break loop
			}
			if t.IsTalkingAboutGo() {
				fmt.Println(t.Username, "\ttweets about golang")
			} else {
				fmt.Println(t.Username, "\tdoes not tweet about golang")
			}
		default:
			break
		}
	}

}

func main() {
	start := time.Now()
	stream := GetMockStream()

	tweets := make(chan *Tweet)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go producer(stream, tweets, &wg)
	wg.Add(1)
	go consumer(tweets, &wg)

	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
