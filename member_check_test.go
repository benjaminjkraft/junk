package junk

import "testing"
import "strconv"

var result bool

const n = 8
const missNum = 1

func _getIntKey(i int) int {
	return i%n + missNum
}

func _getStringKey(i int) string {
	return "aaaaaa" + strconv.Itoa(_getIntKey(i))
}

func BenchmarkMapString(b *testing.B) {
	m := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		m[strconv.Itoa(i)] = struct{}{}
	}
	b.ResetTimer()
	ok := false
	for i := 0; i < b.N; i++ {
		k := _getStringKey(i)
		_, ok = m[k]
	}
	result = ok
}

func BenchmarkSliceString(b *testing.B) {
	m := make([]string, 0, n)
	for i := 0; i < n; i++ {
		m = append(m, strconv.Itoa(i))
	}
	b.ResetTimer()
	ok := false
	for i := 0; i < b.N; i++ {
		k := _getStringKey(i)
		for _, j := range m {
			if k == j {
				ok = true
				break
			}
		}
	}
	result = ok
}

func BenchmarkMapInt(b *testing.B) {
	m := make(map[int]struct{}, n)
	for i := 0; i < n; i++ {
		m[i] = struct{}{}
	}
	b.ResetTimer()
	ok := false
	for i := 0; i < b.N; i++ {
		k := _getIntKey(i)
		_, ok = m[k]
	}
	result = ok
}

func BenchmarkSliceInt(b *testing.B) {
	m := make([]int, 0, n)
	for i := 0; i < n; i++ {
		m = append(m, i)
	}
	b.ResetTimer()
	ok := false
	for i := 0; i < b.N; i++ {
		k := _getIntKey(i)
		for _, j := range m {
			if k == j {
				ok = true
				break
			}
		}
	}
	result = ok
}
