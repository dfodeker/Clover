package errorsx

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// Wrap adds context + caller file:line, while preserving %w for errors.Is/As.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	_, file, line, ok := runtime.Caller(1) // 1 = the caller of Wrap
	loc := "?:0"
	if ok {
		loc = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	return fmt.Errorf("%s (%s): %w", msg, loc, err)
}

// Wrapf variant
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	_, file, line, ok := runtime.Caller(1)
	loc := "?:0"
	if ok {
		loc = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	return fmt.Errorf(fmt.Sprintf("%s (%s): %%w", fmt.Sprintf(format, args...)), loc, err)
}
