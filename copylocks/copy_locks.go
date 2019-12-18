package linters

import (
	"sync"
)

type m = sync.Mutex

func newM() *m { return new(m) }

func caught() {
	var x *m
	_ = *x // [copylocks] assignment copies lock value to _: command-line-arguments.m
	y := new(m)
	_ = *y // [copylocks] assignment copies lock value to _: command-line-arguments.m
	z := newM()
	_ = *z // [copylocks] assignment copies lock value to _: command-line-arguments.m
}

func missed() {
	_ = *new(m)
	_ = *newM()
}
