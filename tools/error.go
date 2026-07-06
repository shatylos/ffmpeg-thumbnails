package tools

import (
	"fmt"
)

type AppError struct {
	Message     string
	ParentError error
	Code        float64
}

func (t AppError) Error() string {
	if t.ParentError != nil && t.ParentError.Error() != "" {
		return fmt.Sprintf("%s. ParentError: %s", t.Message, t.ParentError.Error())
	}
	return t.Message
}
