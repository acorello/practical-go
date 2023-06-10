package main

import (
	"time"
)

/*
For every value "n" in values, spin a goroutine that will
- sleep "n" milliseconds
- Send "n" over a channel

In the function body, collect values from the channel to a slice and return it.
*/
func SleepSort(values []int) (res []int) {
	ch := make(chan int)
	for _, v := range values {
		go /* HERE I SPIN-OFF! */ func(n int) {
			time.Sleep(time.Duration(n) * time.Millisecond)
			// concurrent write to channel means I can't close it right after the for loop: some coroutine may be writing into it.
			ch <- n
		}(v)
	}
	// I can't close the channel in any of the mandating goroutines; because they can't know when they're all done.
	// I can only use a countdown, which in this case could be the length of the slice
	// Can the length of the slice be changed so that id does not match the number of goroutines that were launched?
	for range values {
		v, open := <-ch
		if !open {
			break
		}
		res = append(res, v)
	}
	return res
}

func SleepSortWithSync(values []int) (res []int) {
	ch := make(chan int)
	for _, v := range values {
		go /* HERE I SPIN-OFF! */ func(n int) {
			time.Sleep(time.Duration(n) * time.Millisecond)
			// concurrent write to channel means I can't close it right after the for loop: some coroutine may be writing into it.
			ch <- n
		}(v)
	}
	for range values {
		res = append(res, <-ch)
	}
	return res
}
