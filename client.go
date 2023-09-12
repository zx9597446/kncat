package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"os"
)

func startReverseClient(cfg Config) {
	conn, err := net.Dial(cfg.flgNetwork, cfg.flgConnectAddr)
	if err != nil {
		log.Fatalln(err)
	}
	cconn := NewCryptConn(conn, cfg.flgCryptoMethod, []byte(cfg.flgSecretKey))
	if err = cconn.ClientHandshake(); err != nil {
		log.Fatalf("client handshake error: %s\n", err)
	}
	if cfg.flgCommand == "" {
		log.Fatalf("no command to execute")
	}
	if err := pipeCmd2Conn(cfg.flgCommand, cconn); err != nil {
		logf("pipe local command error: %s", err)
	}
	logf("pipe local command done")
	os.Exit(0)
}

func startClient(cfg Config) {
	conn, err := net.Dial(cfg.flgNetwork, cfg.flgConnectAddr)
	if err != nil {
		log.Fatalln(err)
	}
	cconn := NewCryptConn(conn, cfg.flgCryptoMethod, []byte(cfg.flgSecretKey))
	if err = cconn.ClientHandshake(); err != nil {
		log.Fatalf("client handshake error: %s\n", err)
	}
	if cfg.flgFwdAddr != "" {
		if err := pipeToLocalConn(cfg, cconn); err != nil {
			logf("pipe local connection error: %s", err)
		}
		logf("pipe local connection done")
	} else {
		if err := pipe2LocalConsole(cconn); err != nil {
			logf("pipe local connection error: %s", err)
		}
		logf("pipe local console done")
	}
	os.Exit(0)
}

func pipe2LocalConsole(cconn *CryptConn) (err error) {
	defer cconn.Close()
	sigs := make(chan error, 1)
	go func() {
		err := _fix_pipeStdin(cconn)
		sigs <- err
	}()
	go func() {
		_, err := io.Copy(os.Stdout, cconn)
		sigs <- err
	}()
	err = <-sigs
	return
}

func _fix_pipeStdin(conn *CryptConn) (err error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		logf("Stdin stat failed: %s", err)
		return
	}
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		buffer, err := io.ReadAll(os.Stdin)
		if err != nil {
			logf("Failed read: %s", err)
		}
		io.Copy(conn, bytes.NewReader(buffer))
	} else {
		// Fixed: windows下 os.Stdin没有"\n"导致命令执行失败
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			io.WriteString(conn, input.Text()+"\n")
		}
	}
	return nil
}

func pipeToLocalConn(cfg Config, cconn *CryptConn) (err error) {
	local, err := net.Listen(cfg.flgNetwork, cfg.flgFwdAddr)
	if err != nil {
		log.Fatalln(err)
	}
	lconn, err := local.Accept()
	if err != nil {
		log.Fatal(err)
	}
	return pipe2(cconn, lconn)
}
