package api

import (
	"sync"
	"github.com/labstack/echo"
	"github.com/masaliev/facebook_mini/db"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"strings"
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

	a.echo.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(Key),
		Skipper: func(c echo.Context) bool{
			path := c.Path()
			return strings.Contains(path, "/login") || strings.Contains(path, "/signup")
		},
	}))

	g := a.echo.Group("/api/v1")

	g.POST("/login", a.Login)
	g.POST("/signup", a.SignUp)

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