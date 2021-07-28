package main

import "testing"

var stringResult string

func benchmarkStrings(b *testing.B, length, numAppends int) {
	for i := 0; i < b.N; i++ {
		s := ""
		d := string(make([]byte, length))
		for i := 0; i < numAppends; i++ {
			s += d
		}
		stringResult = s
	}
}

func BenchmarkStrings1(b *testing.B)    { benchmarkStrings(b, 10, 1e1) }
func BenchmarkStrings2(b *testing.B)    { benchmarkStrings(b, 10, 1e2) }
func BenchmarkStrings3(b *testing.B)    { benchmarkStrings(b, 10, 1e3) }
func BenchmarkStrings4(b *testing.B)    { benchmarkStrings(b, 10, 1e4) }
func BenchmarkBigStrings1(b *testing.B) { benchmarkStrings(b, 1000, 1e1) }
func BenchmarkBigStrings2(b *testing.B) { benchmarkStrings(b, 1000, 1e2) }
func BenchmarkBigStrings3(b *testing.B) { benchmarkStrings(b, 1000, 1e3) }
func BenchmarkBigStrings4(b *testing.B) { benchmarkStrings(b, 1000, 1e4) }
