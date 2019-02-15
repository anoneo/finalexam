package main

import (
	"github.com/anoneo/finalexam/customerhandler"
)

func main() {
	customerhandler.CreateTb()
	r := customerhandler.NewRouter()
	r.Run(":2019")
}
