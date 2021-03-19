package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
)

func main() {
	bucket := os.Args[1]
	object := os.Args[2]
	offset, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(err)
	}
	length, err := strconv.Atoi(os.Args[4])
	if err != nil {
		panic(err)
	}
	n, err := strconv.Atoi(os.Args[5])
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	var dt time.Duration
	var lastHash [20]byte
	for i := 0; i < n; i++ {
		s := time.Now()
		rr, err := client.Bucket(bucket).Object(object).NewRangeReader(ctx, int64(offset), int64(length))
		if err != nil {
			panic(err)
		}

		b, err := ioutil.ReadAll(rr)
		if err != nil {
			panic(err)
		}
		dt += time.Since(s)
		h := sha1.Sum(b)
		if h != lastHash {
			fmt.Println("hash: ", hex.EncodeToString(h[:]))
			lastHash = h
		}
	}

	fmt.Println("time/call: ", dt/time.Duration(n))
}
