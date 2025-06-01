package examples

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alesr/errx"
)

func Example() {
	// simulate nested function calls

	deepFunction := func() error {
		return errors.New("file not found")
	}

	middleFunction := func() error {
		return deepFunction()
	}

	anotherFunction := func() error {
		return middleFunction()
	}

	topFunction := func() error {
		if err := anotherFunction(); err != nil {
			return errx.Wrap(err)
		}
		return nil
	}

	if err := topFunction(); err != nil {
		// count frames in verbose output
		verboseStr := fmt.Sprintf("%+v", err)
		frameCount := strings.Count(verboseStr, "[")
		fmt.Printf("Captured %d context frames\n", frameCount)

		// Original error is still accessible
		fmt.Println("Original error preserved:", strings.Contains(err.Error(), "file not found"))
	}

	// Output:
	// Captured 8 context frames
	// Original error preserved: true
}

func Example_chaining() {
	// Simulate nested function calls that each add context
	deepFunction := func() error {
		return errx.Wrap(errors.New("connection refused"))
	}

	middleFunction := func() error {
		if err := deepFunction(); err != nil {
			return errx.Wrap(err)
		}
		return nil
	}

	topFunction := func() error {
		if err := middleFunction(); err != nil {
			return errx.Wrap(err)
		}
		return nil
	}

	err := topFunction()
	if err != nil {
		// Multiple wraps create layered context
		verboseStr := fmt.Sprintf("%+v", err)
		frameCount := strings.Count(verboseStr, "[")
		fmt.Printf("Total context frames captured: %d\n", frameCount)

		// Original error is still accessible
		fmt.Println("Original error preserved:", strings.Contains(err.Error(), "connection refused"))
	}

	// Output:
	// Total context frames captured: 11
	// Original error preserved: true
}
