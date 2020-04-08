package professor

//Professor struct
type Professor struct {
	ID       int    `json:"id" gorm:"primary key"`
	Degree   string `json:"degree" gorm:"type:text"`
	Major    string `json:"major" gorm:"type:text"`
	ImageURL string `json:"image_url" gorm:"type:text"`
	UserID   int    `gorm:"unique_index:idx_prof_user"`
}

//ProfessorResponse struct api
type ProfessorResponse struct {
	ProfessorID int    `json:"id"`
	Degree      string `json:"degree"`
	Major       string `json:"major"`
	ImageURL    string `json:"image_url"`
	User        struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		UserName  string `json:"user_name"`
		Phone     string `json:"phone"`
	}
}

//ProfessorRequest struct api
type ProfessorRequest struct {
	Degree    string `json:"degree" validate:"required"`
	Major     string `json:"major" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	UserName  string `json:"user_name" validate:"required,email"`
	Password  string `json:"password" validate:"required,gte=6,lte=13"`
	Phone     string `json:"phone" validate:"omitempty"`
}

//ProfessorUpdate struct API
type ProfessorUpdate struct {
	Degree    string `json:"degree" validate:"required"`
	Major     string `json:"major"  validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
	UserName  string `json:"user_name"  validate:"required,email"`
	Phone     string `json:"phone"`
}
