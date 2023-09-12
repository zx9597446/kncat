package main

import (
	"io"
)

func pipe2(dst io.ReadWriteCloser, src io.ReadWriteCloser) (err error) {
	defer func() {
		dst.Close()
		src.Close()
	}()
	sigs := make(chan error, 1)
	go func() {
		_, err := io.Copy(dst, src)
		sigs <- err

	}()
	go func() {
		_, err := io.Copy(src, dst)
		sigs <- err
	}()
	err = <-sigs
	return
}
