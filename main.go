package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/paulbellamy/snippets/rest"
	"github.com/paulbellamy/snippets/snippets"
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
	s, err := snippets.NewStore()
	if err != nil {
		return err
	}
	defer s.Close()

	server := rest.NewServer(c.BaseUrl, s)
	log.Println("server listening:", c.BaseUrl)
	return http.ListenAndServe(":3000", server)
}
