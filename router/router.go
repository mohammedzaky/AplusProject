package router

import (
	admin "gitlab.com/mohamedzaky/aplusProject/source/admin"
	auth "gitlab.com/mohamedzaky/aplusProject/source/authentication"
	choice "gitlab.com/mohamedzaky/aplusProject/source/choice"
	handler "gitlab.com/mohamedzaky/aplusProject/source/connectDB"
	course "gitlab.com/mohamedzaky/aplusProject/source/course"
	enroll "gitlab.com/mohamedzaky/aplusProject/source/enroll"
	exam "gitlab.com/mohamedzaky/aplusProject/source/exam"
	migration "gitlab.com/mohamedzaky/aplusProject/source/migration"
	professor "gitlab.com/mohamedzaky/aplusProject/source/professor"
	question "gitlab.com/mohamedzaky/aplusProject/source/question"
	semester "gitlab.com/mohamedzaky/aplusProject/source/semester"
	student "gitlab.com/mohamedzaky/aplusProject/source/student"
	studentAnswer "gitlab.com/mohamedzaky/aplusProject/source/studentAnswer"
	studentExam "gitlab.com/mohamedzaky/aplusProject/source/studentExam"
	user "gitlab.com/mohamedzaky/aplusProject/source/user"

	_ "github.com/lib/pq"

	"github.com/labstack/echo"
)

// New is creating a new object of echo
func New() *echo.Echo {

	e := echo.New()

	//Serve Static files images here
	e.Static("/static", "temp-images")

	//set database constrains
	db := handler.ConnectDB()

	migration.Migration(db)

	//set main routes
	choice.MainGroup(e)
	course.MainGroup(e)
	studentExam.MainGroup(e)
	enroll.MainGroup(e)
	professor.MainGroup(e)
	exam.MainGroup(e)
	question.MainGroup(e)
	studentAnswer.MainGroup(e)
	semester.MainGroup(e)
	student.MainGroup(e)
	user.MainGroup(e)
	auth.MainGroup(e)
	admin.MainGroup(e)
	return e
}
