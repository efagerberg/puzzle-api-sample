package main

import (
	"os"
)

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_USER"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_PORT"),
	)

	a.Run(":8000")
}
