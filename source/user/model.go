package user

//User struct
type User struct {
	ID        int    `json:"id" gorm:"primary key"`
	FirstName string `json:"first_name" validate:"required" gorm:"type:text"`
	LastName  string `json:"last_name" validate:"required" gorm:"type:text"`
	UserName  string `json:"user_name" validate:"required,email" gorm:"unique_index:idx_username;type:text"`
	Password  string `json:"password" validate:"required,gte=6,lte=13" gorm:"type:text"`
	Phone     string `json:"phone" validate:"omitempty"`
}

//UserResponse struct
type UserResponse struct {
	ID        int    `json:"id" gorm:"primary key"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"user_name"`
	Phone     string `json:"phone"`
}

//UserUpdate struct
type UserUpdate struct {
	ID        int    `json:"id" gorm:"primary key"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	UserName  string `json:"user_name" validate:"required,email"`
	Phone     string `json:"phone"`
}

//UpdatePassword struct to change password
type UpdatePassword struct {
	OldPassword     string `json:"old_password" validate:"required,gte=6,lte=13"`
	NewPassword     string `json:"new_password" validate:"required,gte=6,lte=13"`
	ConfirmPassword string `json:"confirm_password" validate:"required,gte=6,lte=13"`
}
