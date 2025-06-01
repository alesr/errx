# errx

[![CI](https://github.com/alesr/errx/actions/workflows/ci.yaml/badge.svg)](https://github.com/alesr/errx/actions/workflows/ci.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/alesr/errx.svg)](https://pkg.go.dev/github.com/alesr/errx)
[![Go Report Card](https://goreportcard.com/badge/github.com/alesr/errx)](https://goreportcard.com/report/github.com/alesr/errx)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## The Problem

You know the drill. Every time you handle an error, you're supposed to add context and wrap the previous one:

```go
func processOrder(id string) error {
    err := validateOrder(id)
    if err != nil {
        return fmt.Errorf("validating order: %w", err)
    }

    err = saveOrder(id)
    if err != nil {
        return fmt.Errorf("saving order: %w", err)
    }

    // ... and so on
}
```

It's good practice, but:

- Repetitive error wrapping everywhere
- Inconsistent context messages
- Missing context when you're in a hurry
- Still no idea where the hell the error actually came from

## The Solution

**errx** does what you should be doing manually, but automatically:

```go
func processOrder(id string) error {
    err := validateOrder(id)
    if err != nil {
        return errx.Wrap(err) // That's it
    }
    return nil
}
```

And you get this:

```
processOrder (orders.go:15): validateOrder (validation.go:23): checkCustomer (customer.go:45): order validation failed: invalid customer ID
```

One `errx.Wrap()` call gives you the full call stack. No manual context needed.

## Why This Matters

Think of it as lightweight tracing for errors. If you're not ready to invest in full distributed tracing (or don't need it), this gives you most of the debugging benefits with zero setup.

Perfect for:
- Small to medium projects
- Microservices that don't have tracing yet
- When you want better errors without the overhead
- Teams that keep forgetting to add proper error context

## Quick Start

```go
package main

import (
    "errors"
    "fmt"
    "github.com/alesr/errx"
)

func deepFunction() error {
    return errors.New("database connection failed")
}

func middleFunction() error {
    return deepFunction() // wrapping is optional until you need it
}

func topFunction() error {
    err := middleFunction()
    if err != nil {
        return errx.Wrap(err) // one wrap captures everything
    }
    return nil
}

func main() {
    if err := topFunction(); err != nil {
        fmt.Println(err)
        // Shows: topFunction (main.go:15)...: middleFunction (main.go:11)...: deepFunction (main.go:7)...: database connection failed
    }
}
```

## Installation

```bash
go get github.com/alesr/errx
```

## How It Works

**First wrap on a new error:**
- Scans up the call stack (up to 10 levels)
- Captures function names, files, line numbers
- Stores it all efficiently

**Subsequent wraps on already-wrapped errors:**
- Just adds the current call to the chain
- No expensive stack scanning

**You get:**
- Function names (cleaned up, no ugly package paths)
- File names and line numbers
- Full compatibility with `errors.Is`, `errors.As`, `errors.Unwrap`

## Examples

### Basic Usage

```go
err := someFunction()
if err != nil {
    return errx.Wrap(err)
}
```

### Multiple Wraps

When functions in the chain each add their own wrap:

```go
func level1() error {
    return errx.Wrap(errors.New("original error"))
}

func level2() error {
    if err := level1(); err != nil {
        return errx.Wrap(err) // adds level2 to the chain
    }
    return nil
}
```

### Verbose Output

Use `%+v` to see each frame on its own line:

```go
fmt.Printf("%+v\n", err)
// Output:
// [0] handleRequest (handler.go:20): API timeout
// [1] processRequest (logic.go:15): API timeout
// [2] callAPI (client.go:10): API timeout
```

### Mix with Manual Context

You can still add manual context when needed:

```go
func processOrder(orderID string) error {
    err := validateOrder(orderID)
    if err != nil {
        contextual := fmt.Errorf("processing order %s: %w", orderID, err)
        return errx.Wrap(contextual) // combines both approaches
    }
    return nil
}
```

## When NOT to Use This

- High-performance hot paths (stack scanning has overhead)
- When you already have comprehensive tracing
- Libraries that others will import (let them decide on error handling)

## Runnable Example

Run `go test -run Example ./examples` to see errx in action:

```go
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
```

## License

MIT
