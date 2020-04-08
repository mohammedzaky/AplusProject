package migration

import (
	"github.com/jinzhu/gorm"
	adminModel "gitlab.com/mohamedzaky/aplusProject/source/admin"
	choiceModel "gitlab.com/mohamedzaky/aplusProject/source/choice"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	courseModel "gitlab.com/mohamedzaky/aplusProject/source/course"
	enrollModel "gitlab.com/mohamedzaky/aplusProject/source/enroll"
	examModel "gitlab.com/mohamedzaky/aplusProject/source/exam"
	professorModel "gitlab.com/mohamedzaky/aplusProject/source/professor"
	questionModel "gitlab.com/mohamedzaky/aplusProject/source/question"
	semesterModel "gitlab.com/mohamedzaky/aplusProject/source/semester"
	studentModel "gitlab.com/mohamedzaky/aplusProject/source/student"
	studentAnswerModel "gitlab.com/mohamedzaky/aplusProject/source/studentAnswer"
	studnetExamModel "gitlab.com/mohamedzaky/aplusProject/source/studentExam"
	tokenModel "gitlab.com/mohamedzaky/aplusProject/source/token"
	userModel "gitlab.com/mohamedzaky/aplusProject/source/user"
)

//Migration to create or remove tables in db
func Migration(db *gorm.DB) {

	dbResult := config.ResetDB()

	resetDB := dbResult

	if resetDB {

		//Remove Constrains for each table
		RemoveConstrains(db)

		db.DropTableIfExists(
			&tokenModel.Token{},
			&adminModel.Admin{},
			&studnetExamModel.StudentExam{},
			&studentAnswerModel.StudentAnswer{},
			&enrollModel.Enrollment{},
			&studentModel.Student{},
			&choiceModel.Choice{},
			&questionModel.Question{},
			&examModel.Exam{},
			&courseModel.Course{},
			&semesterModel.Semester{},
			&professorModel.Professor{},
			&userModel.User{},
		)
	}

	db.AutoMigrate(
		&userModel.User{},
		&tokenModel.Token{},
		&adminModel.Admin{},
		&professorModel.Professor{},
		&studentModel.Student{},
		&courseModel.Course{},
		&semesterModel.Semester{},
		&examModel.Exam{},
		&questionModel.Question{},
		&studentAnswerModel.StudentAnswer{},
		&choiceModel.Choice{},
		&enrollModel.Enrollment{},
		&studnetExamModel.StudentExam{},
	)

	CreateConstrains(db)
}

//CreateConstrains to adding constraions on tables
func CreateConstrains(db *gorm.DB) {

	//User constraions
	db.Model(&userModel.User{}).AddUniqueIndex("idx_username", "username")

	//Token constraions FK
	db.Model(&tokenModel.Token{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&tokenModel.Token{}).AddUniqueIndex("idx_token_value", "token_value")

	//Admin constraions FK
	db.Model(&adminModel.Admin{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&adminModel.Admin{}).AddUniqueIndex("idx_admin_user", "user_id")

	//Professor constraions FK
	db.Model(&professorModel.Professor{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&professorModel.Professor{}).AddUniqueIndex("idx_prof_user", "user_id")

	//course constraions FK
	db.Model(&courseModel.Course{}).AddForeignKey("professor_id", "professors(id)", "SET NULL", "CASCADE")
	db.Model(&courseModel.Course{}).AddForeignKey("semester_id", "semesters(id)", "SET NULL", "CASCADE")
	db.Model(&courseModel.Course{}).AddUniqueIndex("idx_course_name_sid", "name", "semester_id")

	//semester constraions FK
	db.Model(&semesterModel.Semester{}).AddUniqueIndex("idx_course_name_year", "name", "year")

	//student constrain FK
	db.Model(&studentModel.Student{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&studentModel.Student{}).AddUniqueIndex("idx_seat_number", "seat_number")
	db.Model(&studentModel.Student{}).AddUniqueIndex("idx_student_user", "user_id")

	//Exam constrain FK
	db.Model(&examModel.Exam{}).AddForeignKey("course_id", "courses(id)", "CASCADE", "CASCADE")

	//Question constrain FK
	db.Model(&questionModel.Question{}).AddForeignKey("exam_id", "exams(id)", "CASCADE", "CASCADE")

	//Student Answer Constrain FK
	db.Model(&studentAnswerModel.StudentAnswer{}).AddForeignKey("question_id", "questions(id)", "CASCADE", "CASCADE")
	db.Model(&studentAnswerModel.StudentAnswer{}).AddForeignKey("student_id", "students(id)", "CASCADE", "CASCADE")
	db.Model(&studentAnswerModel.StudentAnswer{}).AddForeignKey("student_choice", "choices(id)", "CASCADE", "CASCADE")
	db.Model(&studentAnswerModel.StudentAnswer{}).AddUniqueIndex("question_student", "question_id", "student_id")

	//Choice constrain FK
	db.Model(&choiceModel.Choice{}).AddForeignKey("question_id", "questions(id)", "CASCADE", "CASCADE")

	//Enroll constrain FK
	db.Model(&enrollModel.Enrollment{}).AddForeignKey("course_id", "courses(id)", "CASCADE", "CASCADE")
	db.Model(&enrollModel.Enrollment{}).AddForeignKey("student_id", "students(id)", "CASCADE", "CASCADE")
	db.Model(&enrollModel.Enrollment{}).AddUniqueIndex("course_student", "student_id", "course_id")

	//StudentExam constrain FK
	db.Model(&studnetExamModel.StudentExam{}).AddForeignKey("student_id", "students(id)", "CASCADE", "CASCADE")
	db.Model(&studnetExamModel.StudentExam{}).AddForeignKey("exam_id", "exams(id)", "CASCADE", "CASCADE")

}

//RemoveConstrains to remove constraions on tables
func RemoveConstrains(db *gorm.DB) {

	//remove foreignkey of admin table
	db.Model(&adminModel.Admin{}).RemoveForeignKey("user_id", "users(id)")

	//remove foreignkey of Token table
	db.Model(&tokenModel.Token{}).RemoveForeignKey("user_id", "users(id)")
	db.Model(&tokenModel.Token{}).RemoveIndex("idx_token_value")

	//remove foreignkey of Professor table
	db.Model(&professorModel.Professor{}).RemoveForeignKey("user_id", "users(id)")

	//remove foreignkey of Course table
	db.Model(&courseModel.Course{}).RemoveForeignKey("professor_id", "professors(id)")
	db.Model(&courseModel.Course{}).RemoveForeignKey("semester_id", "semesters(id)")

	//remove foreignkey of student table
	db.Model(&studentModel.Student{}).RemoveForeignKey("user_id", "users(id)")
	db.Model(&studentModel.Student{}).RemoveForeignKey("semester_id", "semesters(id)")
	db.Model(&studentModel.Student{}).RemoveForeignKey("idx_seat_number", "seat_number")

	//remove foreignkey of exam table
	db.Model(&examModel.Exam{}).RemoveForeignKey("course_id", "courses(id)")

	//remove foreignkey of question table
	db.Model(&questionModel.Question{}).RemoveForeignKey("exam_id", "exams(id)")

	//remove many to many foreign keys of studentAnswer table
	db.Model(&studentAnswerModel.StudentAnswer{}).RemoveForeignKey("question_id", "questions(id)")
	db.Model(&studentAnswerModel.StudentAnswer{}).RemoveForeignKey("student_id", "students(id)")
	db.Model(&studentAnswerModel.StudentAnswer{}).RemoveForeignKey("student_choice", "choices(id)")

	//remove foreignkey of Choice table
	db.Model(&choiceModel.Choice{}).RemoveForeignKey("question_id", "questions(id)")

	//remove foreignkey of Enroll table
	db.Model(&enrollModel.Enrollment{}).RemoveIndex("course_student")

	//remove foreignkey of StudentExam table
	db.Model(&studnetExamModel.StudentExam{}).RemoveForeignKey("student_id", "students(id)")
	db.Model(&studnetExamModel.StudentExam{}).RemoveForeignKey("exam_id", "exams(id)")
}
