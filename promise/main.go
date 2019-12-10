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

func race2() {
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	results := make(chan int32)
	defer close(results)
	g.Go(func() error {
		results <- one(ctx)
		return nil
	})
	g.Go(func() error {
		results <- two(ctx)
		return fmt.Errorf("this will cancel the other task")
	})
	var r int32
	select {
	case r = <-results:
		cancel()
	case <-ctx.Done():
		panic("!") // can't happen as-is; in real code we'd return err
	}
	fmt.Println(r)
}

func main() {
	all()
	race()
	race2()
}
