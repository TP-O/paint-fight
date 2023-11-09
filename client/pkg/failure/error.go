package failure

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type AppError struct {
	Code          AppErrorCode
	Msg           string
	OriginalError error
}

const maxTraceback = 10

func (a AppError) Error() string {
	if a.Msg != "" {
		return a.Msg
	} else if a.Code != 0 && errorMsg[a.Code] != "" {
		return errorMsg[a.Code]
	} else {
		return errorMsg[UnknownErrorCode]
	}
}

func (a AppError) ErrorWithTrace() string {
	return fmt.Sprintf("%s \n at %s", a.Error(), trace())
}

func NewTraceError(err error) error {
	return fmt.Errorf("%s \n at %s", err, trace())
}

func trace() string {
	pc := make([]uintptr, maxTraceback)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	// TODO: return more than 1 trace line
	return fmt.Sprintf("\t %s:%d | %s", filepath.Base(frame.File), frame.Line, frame.Function)
}
