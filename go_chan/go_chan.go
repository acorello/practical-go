package main

import (
	"fmt"
)

func main() {
	// go fmt.Println("goroutine")
	// fmt.Println("main")

	// for i := 0; i < 3; i++ {
	// 	// Fix 2: Use a loop body variable
	// 	i := i // "i" shadows "i" from the for loop
	// 	go func() {
	// 		fmt.Println(i) // i from line 14
	// 	}()

	// 	/* Fix 1: Use a parameter
	// 	go func(n int) {
	// 		fmt.Println(n)
	// 	}(i)
	// 	*/
	// 	/* BUG: All goroutines use the "i" for the for loop
	// 	go func() {
	// 		fmt.Println(i) // i from line 12
	// 	}()
	// 	*/
	// }

	// time.Sleep(10 * time.Millisecond)

	// shadowExample()

	// ch := make(chan string)
	// go func() {
	// 	ch <- "hi" // send
	// }()
	// msg := <-ch // receive
	// fmt.Println(msg)

	// go func() {
	// 	for i := 0; i < 3; i++ {
	// 		msg := fmt.Sprintf("message #%d", i+1)
	// 		ch <- msg
	// 	}
	// 	close(ch)
	// }()

	// for msg := range ch {
	// 	fmt.Println("got:", msg)
	// }

	// /* for/range does this
	// for {
	// 	msg, ok := <-ch
	// 	if !ok {
	// 		break
	// 	}
	// 	fmt.Println("got:", msg)
	// }
	// */

	// msg = <-ch // ch is closed
	// fmt.Printf("closed: %#v\n", msg)

	// msg, ok := <-ch // ch is closed
	// fmt.Printf("closed: %#v (ok=%v)\n", msg, ok)

	// ch <- "hi" // ch is closed -> panic
	values := []int{15, 8, 42, 16, 4, 23}
	fmt.Println(SleepSort(values))
}

/* Channel semantics
- send & receive will block until opposite operation (*)
	- Buffered channel has cap(ch) non-blocking sends
- receive from a closed channel will return the zero value without blocking
- send to a closed channel will panic
- closing a closed channel will panic
- send/receive to a nil channel will block forever

See also https://www.353solutions.com/channel-semantics
*/

func shadowExample() {
	n := 7
	{
		n := 2 // from here to } this is "n"
		// n = 2 // outer n
		fmt.Println("inner", n)
	}
	fmt.Println("outer", n)
}
