package main

import (
	"go-enginet/enginet"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Host, Port string
}

func newConfig() *Config {
	conf := new(Config)
	conf.Host = "localhost"
	conf.Port = "8080"
	if os.Getenv("HOST") != "" {
		conf.Host = os.Getenv("HOST")
	}
	if os.Getenv("PORT") != "" {
		conf.Port = os.Getenv("PORT")
	}
	return conf
}

func main() {
	conf := newConfig()
	r := enginet.New()
	r.GET("/", func(ctx *enginet.Context) {
		ctx.String(http.StatusOK, "奥利给")
	})
	r.GET("/hello", func(ctx *enginet.Context) {
		ctx.JSON(http.StatusOK, enginet.H{
			"key": "value",
		})
	})
	log.Fatal(r.Run(conf.Host + ":" + conf.Port))
}
