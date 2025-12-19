package request

type AssignUserToRoleRequest struct {
	UserID uint `json:"user_id" validate:"required"`
	RoleID uint `json:"role_id" validate:"required"`
}

type CreateUserRequest struct {
	Name	 string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Phone	string `json:"phone" validate:"required"`
	Photo   string `json:"photo" validate:"required"`
}

type GetAllUsersRequest struct {
	Page      int    `json:"page" validate:"omitempty,min=1"`
	Limit     int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Search    string `json:"search"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}


type UpdateUserRequest struct {
	Name	 string `json:"name" validate:"required"`
	Email	string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
	Phone	string `json:"phone" validate:"required"`
	Photo   string `json:"photo" validate:"required"`
}