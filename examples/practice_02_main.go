package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// produce is simulating our single input as a channel
func produce() chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- rand.Intn(50)
		}
		fmt.Printf("producer done\n")
		close(ch) // this is important!!!
	}()
	return ch
}

func worker(id int, jobs chan int, out chan OddEven, wg *sync.WaitGroup) {
	for value := range jobs {
		odd := "even"
		if (value & 1) == 1 {
			odd = "odd"
		}
		out <- OddEven{
			Number:  value,
			OddEven: odd,
		}
	}
	close(out) // remember this
	wg.Done()
}

// OddEven struct will be the result of the work done by each fanout thread
// and be the fanin data
type OddEven struct {
	Number  int
	OddEven string
}

func fanin(inputs []chan OddEven) chan OddEven {
	output := make(chan OddEven)

	var wg sync.WaitGroup
	for i, input := range inputs {
		wg.Add(1)
		// explicit params to capture loop vars
		go func(id int, input chan OddEven, output chan OddEven, wg *sync.WaitGroup) {
			for value := range input {
				output <- value
			}
			fmt.Printf("done merging source %d\n", id)
			wg.Done()
		}(i, input, output, &wg)
	}
	go func() {
		wg.Wait()
		close(output) // this is important!!!
	}()

	return output
}

func main() {
	// simulate the input data stream
	inputCh := produce()

	numWorkers := 3
	// fan-out to send data items to workers as individual jobs
	var wg sync.WaitGroup
	workerResults := make([]chan OddEven, numWorkers)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		workerResults[i] = make(chan OddEven)
		go worker(i, inputCh, workerResults[i], &wg)
	}
	go func() {
		wg.Wait()
	}()

	// fan-in the results
	results := fanin(workerResults)

	done := make(chan bool)
	go func() {
		for value := range results {
			fmt.Printf("got %d is %s\n", value.Number, value.OddEven)
		}
		close(done)
	}()
	<-done
	fmt.Println("done")
}
