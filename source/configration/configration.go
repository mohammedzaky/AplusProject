package configration

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

//SignConfig for Signed values like SigningMethod,SigningKey
//for JWT Config ,Token ......etc
type SignConfig struct {
	SignMethod string `yaml:"Method"`
	SignKey    string `yaml:"Key"`
}

//HashConfig for salt value
type HashConfig struct {
	SaltValue string `yaml:"Salt"`
}

//DatabaseConfig struct
type DatabaseConfig struct {
	HostName   string `yaml:"Host"`
	PortNumber string `yaml:"Port"`
	UserDB     string `yaml:"User"`
	NameDB     string `yaml:"Dbname"`
	Password   string `yaml:"Password"`
	SslMode    string `yaml:"Sslmode"`
}

//DataBaseReset struct
type DataBaseReset struct {
	ResetDB bool `yaml:"ResetDB"`
}

//SigningConfig for getting the values (method,key)
func SigningConfig() (string, string) {
	signConfig := SignConfig{}

	yamlFile, err := ioutil.ReadFile("../source/configration/conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &signConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	method := signConfig.SignMethod
	key := signConfig.SignKey

	return method, key
}

//HashingConfig for salt to hash password
func HashingConfig() string {
	hashConfig := HashConfig{}

	yamlFile, err := ioutil.ReadFile("../source/configration/salt.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &hashConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	salt := hashConfig.SaltValue
	return salt
}

//DbConfig for all database connection
func DbConfig() []string {
	databaseConfig := DatabaseConfig{}

	yamlFile, err := ioutil.ReadFile("../source/configration/db.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &databaseConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	host := databaseConfig.HostName
	port := databaseConfig.PortNumber
	user := databaseConfig.UserDB
	dbname := databaseConfig.NameDB
	password := databaseConfig.Password
	mode := databaseConfig.SslMode

	dbResult := make([]string, 0)

	dbResult = append(dbResult, host)
	dbResult = append(dbResult, port)
	dbResult = append(dbResult, user)
	dbResult = append(dbResult, dbname)
	dbResult = append(dbResult, password)
	dbResult = append(dbResult, mode)

	return dbResult
}

// ResetDB function if the value will true then all tables and constrains
// will be dropped and recreate again
func ResetDB() bool {

	resetDB := DataBaseReset{}

	yamlFile, err := ioutil.ReadFile("../source/configration/db.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &resetDB)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	dbRestResult := resetDB.ResetDB

	return dbRestResult
}
