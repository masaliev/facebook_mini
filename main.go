package main

import (
	"github.com/masaliev/facebook_mini/api"
	"os"
)

func main() {
	if _, err := os.Stat("public"); os.IsNotExist(err) {
		os.MkdirAll("public/avatars", 0775)
	}
	a := api.NewApi(":8090", "./facebook_mini.db")
	a.Start()
	a.WaitStop()
}