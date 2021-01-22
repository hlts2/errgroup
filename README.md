# errgroup

Provide synchronization, error propagation, and Context cancelation for groups of goroutines working on subtasks of a common task.
This package highly inspired by [errgroup](https://github.com/golang/sync/tree/master/errgroup).


## Requirement

Go 1.15

## Installing

```
go get github.com/hlts2/errgroup
```

## Example

```go

package main

import (
	"fmt"
	"net/http"

	"github.com/hlts2/errgroup"
)

func main() {
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/",
	}

	eg := new(errgroup.Group)
	for _, url := range urls {
		eg.Go(func() error {
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
    // Perform type conversion and get all goroutine error.
		if err, ok := err.(errgroup.Error); ok {
			fmt.Println(err.Error())
		}
	}
}


```
