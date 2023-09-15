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
		// _, err := io.Copy(cconn, stdin)
		err := _fix_pipeStdin(cconn, stdin)
		sigs <- err
	}()
	err = <-sigs
	logf("pipe stdin/stdout done: %v", err)
	return
}

func _fix_pipeStdin(conn *CryptConn, stdin *os.File) (err error) {
	fi, err := stdin.Stat()
	if err != nil {
		return
	}
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		buffer, err := io.ReadAll(stdin)
		if err != nil {
			return err
		}
		io.Copy(conn, bytes.NewReader(buffer))
	} else {
		// Fixed: windows下 os.Stdin没有"\n"导致命令执行失败
		input := bufio.NewScanner(stdin)
		for input.Scan() {
			io.WriteString(conn, input.Text()+"\n")
		}
	}
	return nil
}
