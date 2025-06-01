# errx

[![Go Reference](https://pkg.go.dev/badge/github.com/alesr/errx.svg)](https://pkg.go.dev/github.com/alesr/errx)
[![Go Report Card](https://goreportcard.com/badge/github.com/alesr/errx)](https://goreportcard.com/report/github.com/alesr/errx)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Why errx?

See the difference instantly. Instead of a basic error message:

```
database connection failed
```

Get comprehensive debugging context automatically:

```
main.main (main.go:12) at 2023-04-15T10:30:47Z: app.initializeApp (app.go:23) at 2023-04-15T10:30:47Z: database.Connect (db.go:45) at 2023-04-15T10:30:47Z: database connection failed
```

**All from a single `errx.Wrap()` call.** No manual context strings needed.

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
    return deepFunction() // No errx needed here
}

func topFunction() error {
    err := middleFunction()
    if err != nil {
        return errx.Wrap(err) // Single wrap captures everything
    }
    return nil
}

func main() {
    if err := topFunction(); err != nil {
        fmt.Println(err)
        // Shows: topFunction (main.go:15) at 2023-...: middleFunction (main.go:11) at 2023-...: deepFunction (main.go:7) at 2023-...: database connection failed
    }
}
```

## Runnable Example

Run `go test -run Example` to see errx in action:

```go
func Example() {
    originalErr := errors.New("database connection failed")
    wrappedErr := Wrap(originalErr)

    fmt.Println("Original error preserved:", errors.Is(wrappedErr, originalErr))
    
    unwrappedErr := errors.Unwrap(wrappedErr)
    fmt.Println("Unwrapped error:", unwrappedErr.Error())
    
    nilErr := Wrap(nil)
    fmt.Println("Wrapping nil returns nil:", nilErr == nil)
    
    // Output:
    // Original error preserved: true
    // Unwrapped error: database connection failed
    // Wrapping nil returns nil: true
}
```

## What is errx?

A lightweight Go package that automatically captures comprehensive debugging context from your call stack. Unlike traditional error wrapping that requires manual context at each level, errx intelligently scans the entire call stack in a single operation.

**Key Benefits:**
- **Automatic context capture** - No manual context strings needed
- **Complete call stack visibility** - See the full error journey  
- **Single wrap coverage** - Works even when intermediate functions don't use errx
- **Standard library compatible** - Works with `errors.Is`, `errors.As`, `errors.Unwrap`
- **Efficient design** - Stack scanning only on first wrap

## Installation

```bash
go get github.com/alesr/errx
```

## Detailed Examples

### Comprehensive Stack Capture

errx automatically captures context from multiple functions, even when they don't use errx:

```go
func apiCall() error {
    return errors.New("API timeout")
}

func processRequest() error {
    // This function doesn't wrap - just passes through
    return apiCall()
}

func handleRequest() error {
    err := processRequest()
    if err != nil {
        // Single wrap captures the entire call chain
        return errx.Wrap(err)
    }
    return nil
}
```

### Verbose Output Format

Use `%+v` for detailed frame-by-frame output:

```go
err := handleRequest()
if err != nil {
    fmt.Printf("%+v\n", err)
}

// Output:
// [0] service.controller (controller.go:20) at 2023-04-15T10:30:47Z: API timeout
// [1] business.processRequest (logic.go:15) at 2023-04-15T10:30:47Z: API timeout
// [2] api.callExternal (client.go:10) at 2023-04-15T10:30:47Z: API timeout
// [3] main.main (main.go:25) at 2023-04-15T10:30:47Z: API timeout
```

### Chaining Multiple Wraps

When you wrap an already wrapped error, errx efficiently adds just the current context:

```go
func level1() error {
    err := errors.New("original error")
    return errx.Wrap(err) // Captures full stack
}

func level2() error {
    err := level1()
    if err != nil {
        return errx.Wrap(err) // Just adds current frame
    }
    return nil
}
```

### Integration with Manual Context

Combine errx with traditional error wrapping when needed:

```go
func processOrder(orderID string) error {
    err := validateOrder(orderID)
    if err != nil {
        // Add manual context first
        contextual := fmt.Errorf("processing order %s: %w", orderID, err)
        // Then capture stack context
        return errx.Wrap(contextual)
    }
    return nil
}
```

## How It Works

**First Wrap (new error):**
1. Captures current caller function context
2. Scans call stack up to 10 levels deep
3. Records function name, file, line number, and timestamp for each frame
4. Stores all context efficiently

**Subsequent Wraps (already wrapped):**
1. Just adds current caller to existing frame chain
2. No expensive stack scanning - maximum efficiency

**Information Captured:**
- Function name (shortened for readability)
- File name (base name only)
- Line number
- Timestamp (RFC3339 format)
- Original error message (preserved)

## Formatting Options

- **Standard format** (`%s`, `%v`): Single line with all context
- **Verbose format** (`%+v`): Multi-line with frame indices
- **Error interface**: Compatible with all standard error handling

## Compatibility

- ✅ Full compatibility with Go's standard error handling
- ✅ Works with `errors.Is`, `errors.As`, `errors.Unwrap`
- ✅ Zero dependencies beyond Go standard library
- ✅ Safe for concurrent use

## License

MIT