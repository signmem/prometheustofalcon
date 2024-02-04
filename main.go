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

	go prom.GetProms()
	// go prom.GetMdsCalAvg()

	g.Logger.Info("program start..")

	if g.Config().SslEnable == true {

		_ , err := os.Stat(g.Config().TLS.CaFile)

		if os.IsNotExist(err) {
			g.Logger.Errorf("%s key not exists", g.Config().TLS.CaFile)
		}

		_ , err = os.Stat(g.Config().TLS.CertFile)

		if os.IsNotExist(err) {
			g.Logger.Errorf("%s key not exists", g.Config().TLS.CertFile)
		}

		_ , err = os.Stat(g.Config().TLS.KeyFile)

		if os.IsNotExist(err) {
			g.Logger.Errorf("%s key not exists", g.Config().TLS.KeyFile)
		}
	}

	select {}
}