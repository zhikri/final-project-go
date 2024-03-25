package main

import (
	"final-project-go/database"
	"final-project-go/routers"
)

func main() {
	database.ConnectDatabase()
	r := routers.SetupRouter()

	r.Run()
}
