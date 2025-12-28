package response

type UserResponse struct {
	ID     string   `json:"id"`
	Email  string   `json:"email"`
	Active bool     `json:"active"`
	Roles  []string `json:"roles"`
}
