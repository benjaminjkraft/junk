package main

// run with go test -bench=.
// My (oversimplified) conclusions are:
// - if you have less than 64 items in each slice, it's faster to do an O(n^2) nested loop
// - if you have more than 64 items, it's worth using a map as an index

import (
	"math/rand"
	"testing"
)

// prevent the compiler from optimizing away the function under test
// by storing its result
var result bool

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func BenchmarkSlices1(b *testing.B) {
	benchmarkSlices(b, 1)
}

func BenchmarkSlices2(b *testing.B) {
	benchmarkSlices(b, 2)
}

func BenchmarkSlices4(b *testing.B) {
	benchmarkSlices(b, 4)
}

func BenchmarkSlices8(b *testing.B) {
	benchmarkSlices(b, 8)
}

func BenchmarkSlices16(b *testing.B) {
	benchmarkSlices(b, 16)
}

func BenchmarkSlices32(b *testing.B) {
	benchmarkSlices(b, 32)
}

func BenchmarkSlices64(b *testing.B) {
	benchmarkSlices(b, 64)
}

func BenchmarkSlices128(b *testing.B) {
	benchmarkSlices(b, 128)
}

func BenchmarkMap1(b *testing.B) {
	benchmarkMap(b, 1)
}

func BenchmarkMap2(b *testing.B) {
	benchmarkMap(b, 2)
}

func BenchmarkMap4(b *testing.B) {
	benchmarkMap(b, 4)
}

func BenchmarkMap8(b *testing.B) {
	benchmarkMap(b, 8)
}

func BenchmarkMap16(b *testing.B) {
	benchmarkMap(b, 16)
}

func BenchmarkMap32(b *testing.B) {
	benchmarkMap(b, 32)
}

func BenchmarkMap64(b *testing.B) {
	benchmarkMap(b, 64)
}

func BenchmarkMap128(b *testing.B) {
	benchmarkMap(b, 128)
}

func benchmarkSlices(b *testing.B, length int) {
	as, bs := genSlices(length)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = IsSubset_UsingSlices(as, bs)
	}
}

func benchmarkMap(b *testing.B, length int) {
	as, bs := genSlices(length)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = IsSubset_UsingMap(as, bs)
	}
}

func genSlices(length int) ([]uint64, []uint64) {
	as := make([]uint64, length)
	bs := make([]uint64, length)
	for i := 0; i < length; i++ {
		as[i] = rand.Uint64()
	}
	for i := 0; i < length; i++ {
		bs[i] = as[length-i-1]
	}
	return as, bs
}

func IsSubset_UsingSlices(as []uint64, bs []uint64) bool {
next_a:
	for _, a := range as {
		for _, b := range bs {
			if b == a {
				continue next_a
			}
		}
		// if we're here, we got through all the Bs without finding a match for A
		return false
	}
	return true
}

func IsSubset_UsingMap(as []uint64, bs []uint64) bool {
	bs_map := make(map[uint64]bool, len(bs))
	for _, b := range bs {
		bs_map[b] = true
	}
	for _, a := range as {
		if _, ok := bs_map[a]; !ok {
			return false
		}
	}
	return true
}
