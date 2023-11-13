package failure

type AppErrorCode = uint

var errorMsg = map[AppErrorCode]string{
	ErrUnknownCode: "Unknown error!",

	ErrIncorrectAccount:                "Email/username or password is incorrect",
	ErrNoLoginResponse:                 "Unable to retrieve login information",
	ErrEmailAlreadyExists:              "Email is unavailable",
	ErrUnableToCreateAccount:           "Unable to create account",
	ErrNoCreateAccountResponse:         "Unable to retrieve account information! Please login to the system",
	ErrInvalidSignature:                "Signature is invalid",
	ErrExpiredLink:                     "Link is expired",
	ErrUnableToVerifyEmail:             "Unable to verify email",
	ErrEmailDoesNotExist:               "Email does not exist",
	ErrUnableToCreateVerifyEmailLink:   "Unable to create verify email link",
	ErrUnableToCreateResetPasswordLink: "Unable to create reset password link",
	ErrUnableToUpdatePassword:          "Unable to update password",
	ErrPlayerDoesNotExist:              "Player does not exist",
	ErrInvalidToken:                    "Token is invalid",
}

/*
Error code details:

- Two first numbers are app code.

- Four last numbers are error code.

Eg: 100000 means 10 (App code) 0000 (Error code)
*/
const (
	ErrUnknownCode AppErrorCode = iota + 10_0000
	ErrIncorrectAccount
	ErrNoLoginResponse
	ErrEmailAlreadyExists
	ErrNoCreateAccountResponse
	ErrInvalidSignature
	ErrExpiredLink
	ErrUnableToVerifyEmail
	ErrEmailDoesNotExist
	ErrUnableToCreateVerifyEmailLink
	ErrUnableToCreateResetPasswordLink
	ErrUnableToUpdatePassword
	ErrPlayerDoesNotExist
	ErrUnableToCreateAccount
	ErrInvalidToken
)
