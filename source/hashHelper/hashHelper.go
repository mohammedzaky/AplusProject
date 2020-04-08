package handlers

import (
	"fmt"
	"log"

	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	"golang.org/x/crypto/scrypt"
)

// DoHash make the hash for each password
func DoHash(password string) (string, error) {
	const (
		PwSaltBytes = 32
		PwHashBytes = 64
	)
	//get salt value from configraion file
	salt := config.HashingConfig()

	result, err := scrypt.Key(
		[]byte(password), //password to be hashing
		[]byte(salt),     //salt value
		16384,            //memory parameter
		8,                //no of iteration (r)
		1,                //no of parallelism (p)
		PwHashBytes,      //size of hash value
	)
	if err != nil {
		log.Fatal(err)
	}
	//casting the value to be string
	hash := fmt.Sprintf("%x\n", result)

	return hash, err

}
