package errx

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const maxDepth = 10

type contextFrame struct {
	funcName string
	file     string
	line     int
	time     time.Time
}

type extendedError struct {
	err    error
	frames []contextFrame
}

// Wrap extends an error by capturing context frames from the call stack.
// It preserves the original error while adding valuable debugging information
// including function names, file locations, line numbers, and timestamps.
// Returns nil if err is nil.
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	// get caller stack info
	// return original error if we can't

	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return err
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return err
	}

	currentFrame := contextFrame{
		funcName: shortenFuncName(fn.Name()),
		file:     filepath.Base(file),
		line:     line,
		time:     time.Now(),
	}

	// check if already wrapped
	// yes: just add the current frame to the existing chain
	// no: capture frames up to 10 levels deep

	if extErr, ok := err.(*extendedError); ok {
		return &extendedError{
			err:    extErr.err,
			frames: append([]contextFrame{currentFrame}, extErr.frames...),
		}
	}

	// first wrap - capture current frame and scan deeper
	frames := []contextFrame{currentFrame}

	// keep going while we can extract valid frame information
	for skip := 2; skip < maxDepth; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn == nil {
			break
		}

		frames = append(frames, contextFrame{
			funcName: shortenFuncName(fn.Name()),
			file:     filepath.Base(file),
			line:     line,
			time:     time.Now(),
		})
	}
	return &extendedError{err: err, frames: frames}
}

// Error returns a string representation of the error with all captured context frames.
// It formats each frame with function name, file location, line number, and timestamp,
// followed by the original error message.
func (e *extendedError) Error() string {
	if len(e.frames) == 0 {
		return e.err.Error()
	}

	var parts []string
	for _, frame := range e.frames {
		part := fmt.Sprintf("%s (%s:%d) at %s",
			frame.funcName,
			frame.file,
			frame.line,
			frame.time.Format(time.RFC3339),
		)
		parts = append(parts, part)
	}
	return fmt.Sprintf("%s: %v", strings.Join(parts, ": "), e.err)
}

// Unwrap returns the original wrapped error, enabling compatibility with errors.Is and errors.As.
func (e *extendedError) Unwrap() error {
	return e.err
}

// Format implements fmt.Formatter to provide detailed error output when using %+v.
// With %+v, it displays each context frame on a separate line with frame indices.
// For other format verbs, it falls back to the standard Error() output.
func (e *extendedError) Format(s fmt.State, verb rune) {
	if verb == 'v' && s.Flag('+') {
		if len(e.frames) == 0 {
			fmt.Fprint(s, e.err.Error())
			return
		}

		for i, frame := range e.frames {
			fmt.Fprintf(s, "[%d] %s (%s:%d) at %s: %v\n",
				i,
				frame.funcName,
				frame.file,
				frame.line,
				frame.time.Format(time.RFC3339),
				e.err,
			)
		}
		return
	}
	fmt.Fprint(s, e.Error())
}

func shortenFuncName(full string) string {
	parts := strings.Split(full, "/")
	if len(parts) > 0 {
		parts = parts[len(parts)-1:]
	}
	return strings.Join(parts, "/")
}
