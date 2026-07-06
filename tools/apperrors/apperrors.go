package apperrors

import (
	"fmt"
	"github.com/shatylos/ffmpeg-thumbnails/tools/terminal"
	"runtime"
)

func New(format string, args ...any) (err error) {
	msg := fmt.Sprintf(format, args...)
	_, file, line, _ := runtime.Caller(1)
	err = fmt.Errorf("%s\n    %s:%d\n", msg, file, line)
	return
}

func Wrap(parent error, format string, args ...any) (err error) {
	msg := fmt.Sprintf(format, args...)
	_, file, line, _ := runtime.Caller(1)
	err = fmt.Errorf("%s\n    %s:%d\n%w", msg, file, line, parent)
	return
}

func NewExcuse(format string, args ...any) (err error) {
	msg := fmt.Sprintf(format, args...)
	_, file, line, _ := runtime.Caller(1)
	err = fmt.Errorf("%s%s%s\n    %s:%d\n", terminal.ColorYellow, msg, terminal.ColorReset, file, line)
	return
}

func WrapExcuse(parent error, format string, args ...any) (err error) {
	msg := fmt.Sprintf(format, args...)
	_, file, line, _ := runtime.Caller(1)
	err = fmt.Errorf("%s%s%s\n    %s:%d\n%w", terminal.ColorYellow, msg, terminal.ColorReset, file, line, parent)
	return
}
