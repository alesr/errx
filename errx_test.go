package errx

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrap(t *testing.T) {
	t.Parallel()

	t.Run("nil input", func(t *testing.T) {
		t.Parallel()

		err := Wrap(nil)
		assert.Nil(t, err)
	})

	t.Run("basic functionality", func(t *testing.T) {
		t.Parallel()

		originalErr := errors.New("test error")
		actualErr := Wrap(originalErr)

		require.NotNil(t, actualErr)
		assert.Equal(t, originalErr, errors.Unwrap(actualErr))

		actualErrStr := actualErr.Error()

		testCases := []struct {
			name     string
			contains string
		}{
			{"function name", "TestWrap"},
			{"file name", "errx_test.go"},
			{"original error", originalErr.Error()},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				assert.Contains(t, actualErrStr, tc.contains)
			})
		}
	})

	t.Run("multiple context frames", func(t *testing.T) {
		t.Parallel()

		// function that calls another function that creates error
		deepFunc := func() error {
			return errors.New("deep error")
		}

		middleFunc := func() error {
			return deepFunc()
		}

		topFunc := func() error {
			if err := middleFunc(); err != nil {
				return Wrap(err)
			}
			return nil
		}

		err := topFunc()
		require.NotNil(t, err)

		errStr := err.Error()

		// should contain function references (the test itself)
		assert.Contains(t, errStr, "TestWrap")

		// when using verbose output, should show multiple frames
		verboseOutput := fmt.Sprintf("%+v", err)

		// should have multiple levels
		lines := strings.Split(strings.TrimSpace(verboseOutput), "\n")
		assert.GreaterOrEqual(t, len(lines), 2)

		// should have level indicators
		assert.Contains(t, verboseOutput, "[0]")
		assert.Contains(t, verboseOutput, "[1]")
	})

	t.Run("error chaining with existing errx errors", func(t *testing.T) {
		t.Parallel()

		baseErr := errors.New("inner error")
		firstLevel := Wrap(baseErr)
		secondLevel := Wrap(firstLevel)
		thirdLevel := Wrap(secondLevel)

		errStr := thirdLevel.Error()

		// should contain multiple function references
		functionCount := strings.Count(errStr, "TestWrap")
		assert.GreaterOrEqual(t, functionCount, 2)

		// should contain original error
		assert.Contains(t, errStr, "inner error")

		// Should contain file references
		fileCount := strings.Count(errStr, "errx_test.go")
		assert.GreaterOrEqual(t, fileCount, 2)
	})

	t.Run("errors.Is compatibility", func(t *testing.T) {
		t.Parallel()

		originalErr := errors.New("sentinel error")
		wrappedErr := Wrap(originalErr)

		assert.True(t, errors.Is(wrappedErr, originalErr))
	})

	t.Run("with manual context", func(t *testing.T) {
		t.Parallel()

		originalErr := errors.New("connection timeout")
		withContext := fmt.Errorf("during database query: %w", originalErr)
		wrappedErr := Wrap(withContext)

		errStr := wrappedErr.Error()

		testCases := []struct {
			name     string
			contains string
		}{
			{"manual context", "during database query"},
			{"original error", "connection timeout"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				assert.Contains(t, errStr, tc.contains)
			})
		}
	})

}

func TestMultipleFrameCapture(t *testing.T) {
	t.Parallel()

	t.Run("captures multiple frames in call chain", func(t *testing.T) {
		t.Parallel()

		// deep call stack

		level4 := func() error {
			return errors.New("source error")
		}

		level3 := func() error {
			return level4()
		}

		level2 := func() error {
			return level3()
		}

		level1 := func() error {
			err := level2()
			if err != nil {
				return Wrap(err)
			}
			return nil
		}

		err := level1()
		require.NotNil(t, err)

		verboseOutput := fmt.Sprintf("%+v", err)

		// have multiple frames (up to 4)
		frameCount := strings.Count(verboseOutput, "[")
		assert.GreaterOrEqual(t, frameCount, 2)
		assert.LessOrEqual(t, frameCount, 4)

		assert.Contains(t, verboseOutput, "TestMultipleFrameCapture")
	})

	t.Run("handles empty call stack gracefully", func(t *testing.T) {
		t.Parallel()

		// this should still work even if we can't capture many frames
		err := Wrap(errors.New("simple error"))

		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "simple error")
	})
}

func TestExtendedErrorFormatting(t *testing.T) {
	t.Parallel()

	t.Run("verbose output format", func(t *testing.T) {
		t.Parallel()

		innerFn := func() error {
			return Wrap(errors.New("inner problem"))
		}

		middleFn := func() error {
			if err := innerFn(); err != nil {
				return Wrap(err)
			}
			return nil
		}

		outerFn := func() error {
			if err := middleFn(); err != nil {
				return Wrap(err)
			}
			return nil
		}

		err := outerFn()
		require.NotNil(t, err)

		verboseOutput := fmt.Sprintf("%+v", err)

		testCases := []struct {
			name  string
			check func(t *testing.T, output string)
		}{
			{
				name: "level indicators",
				check: func(t *testing.T, output string) {
					assert.Contains(t, output, "[0]")
				},
			},
			{
				name: "function names",
				check: func(t *testing.T, output string) {
					assert.Contains(t, output, "TestExtendedErrorFormatting")
				},
			},
			{
				name: "line count",
				check: func(t *testing.T, output string) {
					lineCount := len(strings.Split(output, "\n"))
					assert.GreaterOrEqual(t, lineCount, 2)
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				tc.check(t, verboseOutput)
			})
		}
	})

	t.Run("standard format", func(t *testing.T) {
		t.Parallel()

		err := Wrap(errors.New("test error"))
		standardOutput := fmt.Sprintf("%s", err)

		assert.Contains(t, standardOutput, "test error")
		assert.Contains(t, standardOutput, "TestExtendedErrorFormatting")
	})
}

func TestShortenFuncName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "full package path",
			input:    "github.com/user/repo/pkg.Function",
			expected: "pkg.Function",
		},
		{
			name:     "simple function",
			input:    "main.function",
			expected: "main.function",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "path without dots",
			input:    "a/b/c/d",
			expected: "d",
		},
		{
			name:     "dots without slashes",
			input:    "a.b.c.d",
			expected: "a.b.c.d",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := shortenFuncName(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
