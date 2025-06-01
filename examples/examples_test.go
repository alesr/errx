package examples

import (
	"errors"
	"fmt"

	"github.com/alesr/errx"
)

func Example() {
	originalErr := errors.New("database connection failed")
	wrappedErr := errx.Wrap(originalErr)

	fmt.Println("Original error preserved:", errors.Is(wrappedErr, originalErr))

	unwrappedErr := errors.Unwrap(wrappedErr)
	fmt.Println("Unwrapped error:", unwrappedErr.Error())

	nilErr := errx.Wrap(nil)
	fmt.Println("Wrapping nil returns nil:", nilErr == nil)

	// Output:
	// Original error preserved: true
	// Unwrapped error: database connection failed
	// Wrapping nil returns nil: true
}

func Example_output() {
	// this example shows what errx output looks like
	// Note: actual timestamps and line numbers will vary

	originalErr := errors.New("database connection failed")
	wrappedErr := errx.Wrap(originalErr)

	// the actual output will look similar to this:
	// standard format shows all context in one line:
	fmt.Println("Example standard format:")
	fmt.Println("main.function (file.go:15) at 2023-04-15T10:30:47Z: database connection failed")

	fmt.Println()
	fmt.Println("Example verbose format:")
	fmt.Println("[0] main.function (file.go:15) at 2023-04-15T10:30:47Z: database connection failed")
	fmt.Println("[1] main.caller (file.go:20) at 2023-04-15T10:30:47Z: database connection failed")

	fmt.Println()
	fmt.Println("Error preserved:", errors.Is(wrappedErr, originalErr))

	// Output:
	// Example standard format:
	// main.function (file.go:15) at 2023-04-15T10:30:47Z: database connection failed
	//
	// Example verbose format:
	// [0] main.function (file.go:15) at 2023-04-15T10:30:47Z: database connection failed
	// [1] main.caller (file.go:20) at 2023-04-15T10:30:47Z: database connection failed
	//
	// Error preserved: true
}
