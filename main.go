package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/paulbellamy/snippets/rest"
	"github.com/paulbellamy/snippets/snippets"
)

type Config struct {
	Stdout  io.Writer
	BaseUrl string
	DBUrl   string
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	flag.StringVar(&c.BaseUrl, "base-url", "http://localhost:3000", "base url to reach this service. For returning full urls to the client.")
	flag.StringVar(&c.DBUrl, "db", "memory://", "url to reach the database, either memory:// or redis://user:password@redis-host:6379")
	return
}

func main() {
	c := Config{Stdout: os.Stdout}
	c.RegisterFlags(flag.CommandLine)
	flag.Parse()

	if err := Run(c); err != nil {
		log.Fatal(err)
	}
}

func Run(c Config) error {
	s, err := snippets.NewStore(c.DBUrl)
	if err != nil {
		return err
	}
	defer s.Close()

	server := rest.NewServer(c.Stdout, c.BaseUrl, s)
	log.Println("server listening:", c.BaseUrl)
	return http.ListenAndServe(":3000", server)
}
