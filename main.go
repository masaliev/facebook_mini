package main

import (
	"github.com/masaliev/facebook_mini/api"
	"os"
	"flag"
)

func main() {

	bindAddr := flag.String("bind_addr", ":8090", "Set bind address")
	flag.Parse()

	if _, err := os.Stat("public"); os.IsNotExist(err) {
		os.MkdirAll("public/avatars", 0775)
	}
	a := api.NewApi(*bindAddr, "./facebook_mini.db")
	a.Start()
	a.WaitStop()
}