package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

func one(ctx context.Context) int32 {
	// Simulate a workload.
	time.Sleep(time.Millisecond * time.Duration(rand.Int63n(2000)))
	if ctx.Err() != nil {
		// real workloads this would be a check for context cancellation deep
		// in the http library or whatever; and would do a select to early-exit
		return -1
	}
	return 1
}

func two(ctx context.Context) int32 {
	// Simulate a workload.
	time.Sleep(time.Millisecond * time.Duration(rand.Int63n(1000)))
	time.Sleep(time.Millisecond * time.Duration(rand.Int63n(1000)))
	if ctx.Err() != nil {
		// real workloads this would be a check for context cancellation deep
		// in the http library or whatever; and would do a select to early-exit
		return -1
	}
	return 2
}

func all() {
	g, ctx := errgroup.WithContext(context.Background())
	var oneResult int32
	var twoResult int32
	g.Go(func() error {
		oneResult = one(ctx)
		return nil
	})
	g.Go(func() error {
		twoResult = two(ctx)
		return nil
	})
	g.Wait()
	fmt.Println(oneResult, twoResult)
}

func race() {
	g, ctx := errgroup.WithContext(context.Background())
	var oneResult int32
	var twoResult int32
	g.Go(func() error {
		oneResult = one(ctx)
		return fmt.Errorf("this will cancel the other task")
	})
	g.Go(func() error {
		twoResult = two(ctx)
		return fmt.Errorf("this will cancel the other task")
	})
	g.Wait()
	fmt.Println(oneResult, twoResult)
}

func main() {
	all()
	race()
}
