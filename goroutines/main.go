package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	numbers := []int{1, 2, 3, 4, 5}

	for _, num := range numbers {
		wg.Add(1)
		go squareNum(&wg, num)
	}

	wg.Wait()
	fmt.Println("All go routines are finished")
}

func squareNum(wg *sync.WaitGroup, num int) {
	defer wg.Done()

	square := num * num
	fmt.Printf("%d is squared as %d \n", num, square)
}
