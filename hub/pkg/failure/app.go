package failure

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type AppErrorCode = uint

var errorMsg = map[AppErrorCode]string{
	ErrUnknownCode: "Unknown error!",

	ErrInvalidToken: "Token is invalid",

	ErrPlayerDoesNotExist:    "Player does not exist",
	ErrUnableToCreatePlayer:  "Unable to create player",
	ErrPlayerAlreadyExists:   "Player already exists",
	ErrUsernameAlreadyExists: "Username is not available",
}

/*
Error code details:

- Two first numbers are app code.

- Four last numbers are error code.

Eg: 100000 means 10 (App code) 0000 (Error code)
*/
const (
	ErrUnknownCode AppErrorCode = iota + 10_0000

	ErrInvalidToken

	ErrPlayerDoesNotExist
	ErrUnableToCreatePlayer
	ErrPlayerAlreadyExists
	ErrUsernameAlreadyExists
)

type AppError struct {
	Code          AppErrorCode
	HttpStatus    int
	Msg           string
	OriginalError error
}

const maxTraceback = 10

func (a *AppError) Error() string {
	if a.Msg != "" {
		return a.Msg
	} else if a.Code != 0 && errorMsg[a.Code] != "" {
		return errorMsg[a.Code]
	} else {
		return errorMsg[ErrUnknownCode]
	}
}

func ErrorWithTrace(err error) error {
	fmt.Printf("%s \n at %s", err, trace())
	return err
}

func trace() string {
	pc := make([]uintptr, maxTraceback)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Sprintf("\t %s:%d | %s\n", filepath.Base(frame.File), frame.Line, frame.Function)
}
