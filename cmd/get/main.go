package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	req, err := http.NewRequest("GET", os.Args[1], nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Range", os.Args[2])
	req.Header.Set("Accept-Encoding", "")
	n, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(err)
	}
	var dt time.Duration
	var lastHash [20]byte
	for i := 0; i < n; i++ {
		s := time.Now()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		b, err := ioutil.ReadAll(resp.Body)
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
