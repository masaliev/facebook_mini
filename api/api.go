package api

import (
	"sync"
	"github.com/labstack/echo"
	"github.com/masaliev/facebook_mini/db"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

type Api struct {
	dataStorage *db.DataStorage
	waitGroup sync.WaitGroup
	echo *echo.Echo
	bindAddress string
}

const Key string = "51xlcpBtCQ"


func NewApi(bindAddress string, dbPath string) *Api {
	a := &Api{}
	a.dataStorage = db.NewDataStorage(dbPath)
	a.echo = echo.New()
	a.echo.Logger.SetLevel(log.ERROR)
	a.echo.Use(middleware.Logger())

	a.bindAddress = bindAddress
	return a
}

func (a *Api) WaitStop()  {
	a.waitGroup.Wait()
}

func (a *Api) Start()  {
	a.waitGroup.Add(1)
	go func() {
		a.echo.Start(a.bindAddress)
		a.waitGroup.Done()
	}()
}