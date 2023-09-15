package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

func pipeStdInOut(stdin, stdout *os.File, cconn *CryptConn) (err error) {
	defer cconn.Close()
	logf("pipe stdin/stdout: -> %s", cconn.RemoteAddr().String())
	sigs := make(chan error, 1)
	go func() {
		_, err := io.Copy(stdout, cconn)
		sigs <- err
	}()
	go func() {
		_, err := io.Copy(cconn, stdin)
		sigs <- err
	}()
	err = <-sigs
	logf("pipe stdin/stdout done: %v", err)
	return
}

func pipeOut(out *os.File, cconn *CryptConn) (err error) {
	defer cconn.Close()
	sigs := make(chan error, 1)
	go func() {
		_, err := io.Copy(out, cconn)
		sigs <- err
	}()
	err = <-sigs
	return
}

func pipeIn(in *os.File, cconn *CryptConn) (err error) {
	defer cconn.Close()
	sigs := make(chan error, 1)
	go func() {
		_, err := io.Copy(cconn, in)
		sigs <- err
	}()
	err = <-sigs
	return
}

func _fix_pipeStdin(conn *CryptConn) (err error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return
	}
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		buffer, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
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
