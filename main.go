package main

import "github.com/masaliev/facebook_mini/api"

func main() {
	a := api.NewApi(":8090", "./facebook_mini.db")
	a.Start()
	a.WaitStop()
}