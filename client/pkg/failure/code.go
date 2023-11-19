package failure

type AppErrorCode = uint

var errorMsg = map[AppErrorCode]string{
	ErrUnknownCode: "Unknown error!",

	ErrInvalidToken: "Token is invalid",

	ErrPlayerDoesNotExist:   "Player does not exist",
	ErrUnableToCreatePlayer: "Unable to create player",
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
)
