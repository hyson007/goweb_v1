package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var config struct {
	port int
	host string
}

const (
	usage = `usage: %s
Run HTTP server

Options:
`
)

func PortVar(port *int) *portVar {
	return &portVar{port}
}

type portVar struct {
	port *int
}

func (p *portVar) String() string {
	if p.port == nil {
		return ""
	}

	return fmt.Sprintf("%d", *p.port)
}

func (p *portVar) Set(s string) error {
	val, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	const minPort, maxPort = 1, 65535
	if val < minPort || val > maxPort {
		return fmt.Errorf("port %d out of range [%d:%d]", val, minPort, maxPort)
	}

	*p.port = val
	return nil
}

func main() {
	flag.Var(PortVar(&config.port), "port", "port to listening on")
	//flag.IntVar(&config.port, "port", config.port, "port to listen on")
	flag.StringVar(&config.host, "host", config.host, "host to listen on")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}

func init() { // [6]
	// Set defaults
	s := os.Getenv("HTTPD_PORT")
	p, err := strconv.Atoi(s)
	if err == nil {
		config.port = p
	} else {
		config.port = 8080
	}

	h := os.Getenv("HTTPD_HOST")
	if len(h) > 0 {
		config.host = h
	} else {
		config.host = "localhost"
	}
}
