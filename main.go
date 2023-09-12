package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	flgListenAddr, flgConnectAddr, flgNetwork string
	flgSecretKey, flgCryptoMethod             string
	flgVerbose                                bool
	flgFwdAddr, flgCommand                    string
}

var cfg = Config{}

var (
	logger = log.New(os.Stderr, "[verbose]: ", log.LstdFlags|log.Lshortfile)
)

func logf(f string, v ...interface{}) {
	if cfg.flgVerbose {
		logger.Output(2, fmt.Sprintf(f, v...))
	}
}

func init() {
	flag.StringVar(&cfg.flgListenAddr, "l", ":9597", "listen address")
	flag.StringVar(&cfg.flgConnectAddr, "c", "", "connect address")
	flag.StringVar(&cfg.flgSecretKey, "s", "", "secret key")
	flag.StringVar(&cfg.flgCryptoMethod, "m", "rc4", "crypto method (rc4 or aes256cfb)")
	flag.StringVar(&cfg.flgNetwork, "n", "tcp", "network protocol: tcp tcp4 tcp6")
	flag.StringVar(&cfg.flgCommand, "e", "", "program to execute (such as cmd.exe or /bin/bash)")
	flag.StringVar(&cfg.flgFwdAddr, "f", "", "forward address")

	flag.BoolVar(&cfg.flgVerbose, "v", false, "verbose output")

	flag.Parse()
}

func main() {
	if cfg.flgConnectAddr != "" {
		logf("connect to %s", cfg.flgConnectAddr)
		go startClient(cfg)
	} else if cfg.flgListenAddr != "" {
		logf("listen on %s", cfg.flgListenAddr)
		go startServer(cfg)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigs
}
