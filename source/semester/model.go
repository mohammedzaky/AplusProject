package semester

//Semester struct
type Semester struct {
	ID   int    `json:"id" gorm:"primary key"`
	Name string `json:"name" validate:"required" gorm:"unique_index:idx_course_name_year"`
	Year int    `json:"year" validate:"required,gte=2019,lte=2030" gorm:"unique_index:idx_course_name_year"`
}
type FullName struct {
	firstName string `json:"first_name"`
	lastName  string `json:"last_name"`
}
