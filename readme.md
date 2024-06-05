# Golang worker pool


```go
package main

import(
	"github.com/alserov/wpool"
	"fmt"
	"time"
)

func main() {
	pool := wpool.NewPool(3) // inits new pool with 3 workers
	
	go func() {
		time.Sleep(time.Second * 3)
		pool.Stop() // stops all workers
        }()
	
	for i := 0; i < 10; i++ {
		// adds new task to pool
		pool.Execute(func() error {
			// your function 
		    return nil
		})
        }
	
	// awaits errors from pool
	for err := range pool.AwaitError() {
		fmt.Println(err)
        }
}
```
