package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang/sync/errgroup"
)

type Result struct{}

// BEGIN ANALYZEURL FUNC OMIT

func analyzeURL(ctx context.Context, url string) (*Result, error) {

	// Construct request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Run request
	resp, err := http.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TODO: read resp.Body and analyze

	return &Result{}, nil

}

// END ANALYZEURL FUNC OMIT

// BEGIN ANALYZE FUNC OMIT

func analyze(ctx context.Context) error {

	g, ctx := errgroup.WithContext(ctx)

	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.example.com/",
	}

	for _, url := range urls {
		url := url
		g.Go(func() error {
			_, err := analyzeURL(ctx, url)
			return err
		})
	}

	return g.Wait()

}

// END ANALYZE FUNC OMIT

func parallel() {

	// BEGIN PARALLEL OMIT

	var g errgroup.Group
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.example.com/",
	}
	for _, url := range urls {
		// Launch a goroutine to fetch the URL.
		url := url
		g.Go(func() error {
			// Fetch the URL.
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}
			return err
		})
	}
	// Wait for all HTTP fetches to complete.
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully fetched all URLs.")
	}

	// END PARALLEL OMIT

}

func parallel() {

}
