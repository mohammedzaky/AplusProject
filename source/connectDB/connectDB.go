package source

import (
	"log"

	"github.com/jinzhu/gorm"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
)

//ErrorModel Handle the error coming from database
type ErrorModel struct {
	Message string `json:"Message"`
	Detail  string `json:"Detail"`
}

//ConnectDB to connect to database
func ConnectDB() (db *gorm.DB) {

	dbResult := config.DbConfig()
	dbString := "host=" + dbResult[0] + " port=" + dbResult[1] + " user=" + dbResult[2] + " dbname=" + dbResult[3] + " password=" + dbResult[4] + " sslmode=" + dbResult[5]
	db, err := gorm.Open("postgres", dbString)
	if err != nil {
		log.Panic(err)
	}
	return db
}
