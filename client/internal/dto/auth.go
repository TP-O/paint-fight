package dto

type Login struct {
	Email    string `json:"email" validate:"required_without=Username,email,max=320"`
	Username string `json:"username" validate:"required_without=Email,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateAccount struct {
	Email                string `json:"email" validate:"required,email,max=320"`
	Password             string `json:"password" validate:"required,min=8"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"eqfield=Password"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required,email,max=320"`
}

type ResetPassword struct {
	Password             string `json:"password" validate:"required,min=8"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"eqfield=Password"`
}
