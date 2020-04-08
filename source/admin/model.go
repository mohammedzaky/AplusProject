package admin

//Admin struct
type Admin struct {
	ID       int    `json:"id" gorm:"primary key"`
	Position string `json:"position" gorm:"type:text"`
	UserID   int    `gorm:"unique_index:idx_admin_user"`
}

//AdminRequest struct api
type AdminRequest struct {
	ID        int    `json:"id"`
	Position  string `json:"position"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	UserName  string `json:"user_name" validate:"required,email"`
	Password  string `json:"password" validate:"required,gte=6,lte=13"`
	Phone     string `json:"phone"`
}

//AdminUpdate struct API
type AdminUpdate struct {
	Position  string `json:"position"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	UserName  string `json:"user_name" validate:"required,email"`
	Phone     string `json:"phone"`
}

//AdminResponse struct
type AdminResponse struct {
	AdminID  int    `json:"id"`
	Position string `json:"position"`

	User struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"user_name"`
		Phone     string `json:"phone"`
	}
}
