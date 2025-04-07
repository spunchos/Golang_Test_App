package user

type UserDTO struct {
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Email        string `json:"email"`
}
