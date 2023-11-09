package failure

type AppErrorCode = uint

var errorMsg = map[AppErrorCode]string{
	UnknownErrorCode: "Unknown error!",
}

/*
Error code details:

- Two first numbers are service code.

- Two middle numbers are function code.

- Two last numbers are error code.

Eg: 100000 means 10 (Service code) 00 (Function code) 00 (Error code)
*/
const (
	UnknownErrorCode AppErrorCode = 10_00_00
)
