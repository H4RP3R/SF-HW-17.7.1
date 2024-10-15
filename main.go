// Напишите код, в котором несколько горутин увеличивают значение целочисленного
// счётчика и синхронизируют свою работу через канал. Нужно предусмотреть
// возможность настройки количества используемых горутин и конечного значения
// счётчика, до которого его следует увеличивать.

// Попробуйте реализовать счётчик с элементами ООП (в виде структуры и методов структуры).
// Попробуйте реализовать динамическую проверку достижения счётчиком нужного значения.

package main

import (
	"counter"
	"flag"
	"fmt"
	"sync"
)

var (
	goroutineNum  int
	maxCounterVal int

	wg sync.WaitGroup
)

func init() {
	flag.IntVar(&goroutineNum, "gnum", 1, "number of goroutines")
	flag.IntVar(&maxCounterVal, "max", 1_000_000, "max counter value")
	flag.Parse()
}

func main() {
	cnt := counter.New(maxCounterVal, 1)

	wg.Add(goroutineNum)
	for i := 0; i < goroutineNum; i++ {
		go func() {
			defer wg.Done()
			for {
				if current := cnt.Increment(); current < maxCounterVal {
					cnt.Increment()
				} else {
					return
				}
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Counter: %d\n", cnt.Val())
	cnt.Close()
}
