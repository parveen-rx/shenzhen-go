// The transform command was automatically generated by Shenzhen Go.
package main

import (
	"fmt"
	"runtime"
	"sync"
)

var _ = runtime.Compiler

func Generate_some_numbers(nums chan<- int) {
	// Generate some numbers

	defer func() {
		close(nums)
	}()
	for i := 0; i < 10; i++ {
		nums <- i
	}
}

func Multiply_numbers_by_3(inputs <-chan int, outputs chan<- int) {
	// Multiply numbers by 3

	defer func() {
		if outputs != nil {
			close(outputs)
		}
	}()
	for input := range inputs {
		out := func() int {
			return input * 3
		}()
		if outputs != nil {
			outputs <- out
		}
	}
}

func Print_numbers(inputs <-chan int, outputs chan<- interface{}) {
	// Print numbers

	defer func() {
		if outputs != nil {
			close(outputs)
		}
	}()
	for input := range inputs {
		out := func() interface{} {
			fmt.Println(input)
			return nil
		}()
		if outputs != nil {
			outputs <- out
		}
	}
}

func main() {

	channel0 := make(chan int, 0)
	channel1 := make(chan int, 0)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		Generate_some_numbers(channel0)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		Multiply_numbers_by_3(channel0, channel1)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		Print_numbers(channel1, nil)
		wg.Done()
	}()

	// Wait for the various goroutines to finish.
	wg.Wait()
}
