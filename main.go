package main

import (
	"fmt"
	"go-enginet/enginet"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Host, Port string
}

func NewConfig() *Config {
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
	conf := NewConfig()
	r := enginet.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k,v := range r.Header {
			_, _ = fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})
	log.Fatal(r.Run(conf.Host + ":" + conf.Port))
}

