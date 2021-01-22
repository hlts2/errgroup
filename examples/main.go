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
		if err, ok := err.(errgroup.Error); ok {
			fmt.Println(err.Error())
		}
	}
}
