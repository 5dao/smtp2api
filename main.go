package main

import (
	"flag"
	"log"
	"os"

	"github.com/5dao/smtp2api/api"
	"github.com/5dao/smtp2api/utils"
)

var cfg *api.Config
var cfgFile string

var makeKey bool

func init() {
	cfg = &api.Config{}

	flag.StringVar(&cfgFile, "c", "conf.toml", "-c conf.toml")
	flag.BoolVar(&makeKey, "key", false, "--key make random key")
	flag.Parse()
}

func main() {
	if makeKey {
		log.Println(utils.AesRandomKey())
		os.Exit(0)
		return
	}

	if err := utils.LoadConfig(cfgFile, cfg); err != nil {
		log.Println(err)
		os.Exit(0)
		return
	}

	svr, err := api.NewServer(cfg)
	if err != nil {
		log.Println("NewServer err: ", err)
		return
	}
	svr.Start()

	select {}
}
