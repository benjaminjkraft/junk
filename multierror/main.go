package main

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

func main() {
	inner := fmt.Errorf("inner err")        // errors.Internal("my sentinel")
	outer := fmt.Errorf("wraps: %w", inner) // errors.Wrap(inner)
	multi := multierror.Append(outer, outer)
	fmt.Println(errors.Is(multi, inner)) // returns true!
}
