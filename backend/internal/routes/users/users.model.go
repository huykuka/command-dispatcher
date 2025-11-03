package users

type CreateUserDTO struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Age   int    `json:"age" validate:"required,min=3,max=150"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserDTO struct {
	Name string `json:"name,omitempty" validate:"omitempty,min=3,max=32"`
	Age  int    `json:"age,omitempty" validate:"omitempty,min=3,max=32"`
}

type GetUserQuery struct {
	Page struct {
		Number int `json:"number,omitempty" form:"page[number]" validate:"omitempty,min=5"`
		Size   int `json:"size,omitempty" form:"page[size]" validate:"omitempty,min=5"`
	} `json:"page,omitempty"`
	Sort struct {
		Key   string `json:"key,omitempty" form:"sort[order]" validate:"omitempty"`
		Value string `json:"value,omitempty" form:"sort[order]" validate:"omitempty,lowercase"`
	} `json:"sort,omitempty"`
	Filter struct {
		Role string `json:"role,omitempty" form:"filter[vendorID]" validate:"omitempty"`
	} `json:"filter,omitempty"`
}
