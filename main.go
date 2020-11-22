package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/paulbellamy/snippets/rest"
)

type Config struct {
	BaseUrl string
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	flag.StringVar(&c.BaseUrl, "base-url", "http://localhost:3000", "base url to reach this service. For returning full urls to the client.")
	return
}

func main() {
	c := Config{}
	c.RegisterFlags(flag.CommandLine)
	flag.Parse()

	if err := Run(c); err != nil {
		log.Fatal(err)
	}
}

func Run(c Config) error {
	server := rest.NewServer(c.BaseUrl)
	log.Println("server listening:", c.BaseUrl)
	return http.ListenAndServe(":3000", server)
}
