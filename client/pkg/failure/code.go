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

	ErrPlayerDoesNotExist: "Player does not exist",
}

/*
Error code details:

- Two first numbers are app code.

- Two middle numbers are app service code.

- Three last numbers are error code.

Eg: 100000 means 10 (App code) 00 (App service code) 00 (Error code)
*/
const (
	ErrUnknownCode AppErrorCode = 10_00_000

	ErrIncorrectAccount                AppErrorCode = 10_01_000
	ErrNoLoginResponse                 AppErrorCode = 10_01_001
	ErrEmailAlreadyExists              AppErrorCode = 10_01_002
	ErrUnableToCreateAccount           AppErrorCode = 10_01_003
	ErrNoCreateAccountResponse         AppErrorCode = 10_01_004
	ErrInvalidSignature                AppErrorCode = 10_01_005
	ErrExpiredLink                     AppErrorCode = 10_01_006
	ErrUnableToVerifyEmail             AppErrorCode = 10_01_007
	ErrEmailDoesNotExist               AppErrorCode = 10_01_008
	ErrUnableToCreateVerifyEmailLink   AppErrorCode = 10_01_009
	ErrUnableToCreateResetPasswordLink AppErrorCode = 10_01_010
	ErrUnableToUpdatePassword          AppErrorCode = 10_01_011

	ErrPlayerDoesNotExist AppErrorCode = 10_02_000
)
