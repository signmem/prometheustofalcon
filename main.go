package main

import (
	"flag"
	"fmt"
	"github.com/signmem/prometheustofalcon/g"
	"github.com/signmem/prometheustofalcon/prom"
	"os"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")

	flag.Parse()

	if *version {
		version := g.Version
		fmt.Printf("%s", version)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	g.Logger = g.InitLog()

	prom.GetProms()
	select {}
}