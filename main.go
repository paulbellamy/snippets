package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/paulbellamy/snippets/rest"
)

type Config struct {
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
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
	server := rest.NewServer()
	log.Println("server listening:", "http://localhost:3000")
	return http.ListenAndServe(":3000", server)
}
