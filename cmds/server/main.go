package main

import (
	"flag"
	"os"

	"github.com/5dao/golibs/log"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"

	"github.com/5dao/smtp2api/server"
)

var cfg *server.Config
var cfgFile string

var makeKey bool

func init() {
	cfg = &server.Config{}

	flag.StringVar(&cfgFile, "c", "conf.toml", "-c conf.toml")
	flag.BoolVar(&makeKey, "key", false, "--key make random key")
	flag.Parse()
}

func main() {
	if makeKey {
		log.Println(server.AesRandomKey())
		os.Exit(0)
		return
	}

	if _, err := toml.DecodeFile(cfgFile, cfg); err != nil {
		log.Println(err)
		os.Exit(0)
		return
	}

	gin.SetMode(gin.ReleaseMode)

	svr, err := server.NewServer(cfg)
	if err != nil {
		log.Println("NewServer err: ", err)
		return
	}
	svr.Start()

	select {}
}
