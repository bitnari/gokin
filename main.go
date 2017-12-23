package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"encoding/json"
	"io/ioutil"
	"path"
	"os"
	"io"
	"path/filepath"
	"log"
	"fmt"
	"net/http"
	"errors"
)

type Config struct {
	Server struct {
		Port            int `json:"port"`
		CreditPriority  string `json:"creditPriority"`
		Games           []string `json:"games"`
	} `json:"server"`
	DB struct {
		Host        string `json:"host"`
		DB          string `json:"db"`
	} `json:"database"`
}
var config Config

func (c *Config) HasGame(game string) bool {
	for _, g := range c.Server.Games {
		if g == game {
			return true
		}
	}
	return false
}

var mongo MongoConnection

func copyFile(path, dest string) error {
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	e := out.Close()
	if e != nil {
		return e
	}

	return nil
}

var AppPath string
func init() {
	var err error
	AppPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panic(err)
	}

	if _, err = os.Stat(path.Join(AppPath, "config.json")); os.IsNotExist(err) {
		copyFile(path.Join(AppPath, "resources", "config.json"), path.Join(AppPath, "config.json"))
	}

	b, err := ioutil.ReadFile(path.Join(AppPath, "config.json"))
	if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Panic(err)
	}
}

/* func CORSKakin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		res := c.Response()
		res.Header().Set(echo.HeaderAccessControlAllowOrigin, "kakin://kakin")
		
		return next(c)
	}
	
	Custom protocols are not supported. Should use wildcard.
} */
	
func main() {
	if config.Server.CreditPriority != "gold" && config.Server.CreditPriority != "credit" {
		panic(errors.New(`config creditPriority must be "credit" or "gold"`))
	}

	mongo = MongoConnection{}
	err := mongo.Init(config.DB.Host, config.DB.DB)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))

	fs := http.FileServer(http.Dir(path.Join(AppPath, "kakin", "dist")))
	e.GET("/dist/*", echo.WrapHandler(http.StripPrefix("/dist/", fs)))
	e.File("/", path.Join(AppPath, "kakin", "index.html"))

	e.POST("/api/v1/verify", VerifyAccount)
	e.POST("/api/v1/account", GetAccount)
	e.POST("/api/v1/renew", RenewToken)
	e.POST("/api/v1/usegold", SubtractGold)

	e.POST("/api/v2/register", RegisterHandler)
	e.POST("/api/v2/verify", VerifyAccount_v2)
	e.POST("/api/v2/pay", Pay_v2)
	e.POST("/api/v2/renew", RenewToken)
	e.POST("/api/v2/account", Account_v2)
	e.POST("/api/v2/score", Score_v2)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Server.Port)))
}
