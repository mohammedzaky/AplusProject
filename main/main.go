package main

import (
	"log"

	router "gitlab.com/mohamedzaky/aplusProject/router"
)

func main() {

	e := router.New()
	log.Fatal(e.Start(":1323"))

}
