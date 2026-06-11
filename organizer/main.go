package main

import (
	"fmt"
	"time"
)

func main() {

	c := make(chan int)

	go func() { 
		c <- 2;
		close(c)
 	}()
	
	time.Sleep(1 * time.Second)
	v, ok := <-c
	v1, ok1 := <-c
	v2, ok2 := <-c
	fmt.Println(v, ok)
	fmt.Println(v1, ok1)
	fmt.Println(v2, ok2)
	
}
