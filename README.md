# errgroup

[![Go Report Card](https://goreportcard.com/badge/github.com/hlts2/errgroup)](https://goreportcard.com/report/github.com/hlts2/errgroup)
[![GoDoc](http://godoc.org/github.com/hlts2/errgroup?status.svg)](http://godoc.org/github.com/hlts2/errgroup)

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
		url := url
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
		if err, ok := err.(errgroup.Error); ok {
			fmt.Println(err.Errors()) // slice of errors that occurred inside eg.Go
		}
	} else {
		fmt.Println("Successfully fetched all URLs.")
	}
}

```

## Contribution
1. Fork it ( https://github.com/hlts2/errgroup/fork )
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request

## Author
[hlts2](https://github.com/hlts2)
