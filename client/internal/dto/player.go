package dto

type CreatePlayer struct {
	Username string `json:"username" binding:"required,alphanum,min=1,max=25"`
}
