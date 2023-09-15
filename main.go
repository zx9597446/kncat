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
	flgReverse                                bool
}

var (
	logger = log.New(os.Stderr, "[verbose]: ", log.LstdFlags|log.Lshortfile)
	cfg    = Config{}
	chSigs = make(chan os.Signal, 1)
)

func logf(f string, v ...interface{}) {
	if cfg.flgVerbose {
		logger.Output(2, fmt.Sprintf(f, v...))
	}
}

func initFlags() {
	flag.StringVar(&cfg.flgListenAddr, "l", ":9597", "for server: listening on")
	flag.StringVar(&cfg.flgConnectAddr, "c", "", "for client: connect to")
	flag.StringVar(&cfg.flgSecretKey, "s", "", "your secret key")
	flag.StringVar(&cfg.flgCryptoMethod, "m", "aes256cfb", "crypto method (rc4 or aes256cfb)")
	flag.StringVar(&cfg.flgNetwork, "n", "tcp", "network protocol: tcp tcp4 tcp6")
	flag.StringVar(&cfg.flgCommand, "e", "", "for server: program to execute (cmd.exe or /bin/bash or with args: cat -- some.log, use -- split args)")
	flag.StringVar(&cfg.flgFwdAddr, "f", "", "forward address(server: connect to this address. client: accept on this address)")

	flag.BoolVar(&cfg.flgReverse, "r", false, "reverse mode: connect and execute program on client side, to get a reverse shell")
	flag.BoolVar(&cfg.flgVerbose, "v", false, "verbose output")

	flag.Parse()
}

func main() {
	initFlags()

	if cfg.flgConnectAddr != "" {
		logf("connect to %s", cfg.flgConnectAddr)
		go startClient(cfg)
	} else if cfg.flgListenAddr != "" {
		logf("listen on %s", cfg.flgListenAddr)
		go startServer(cfg)
	}
	waitSignals()
}

func waitSignals() {
	signal.Notify(chSigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	s := <-chSigs
	logf("signal: %v", s)
}
