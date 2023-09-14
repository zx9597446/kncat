package main

import (
	"os/exec"
	"time"
)

func pipeCmd2Conn(strCmd string, cconn *CryptConn) (err error) {
	defer cconn.Close()
	exe, args := parseCmd(strCmd)
	cmd := exec.Command(exe, args...)
	cmd.Stdin = cconn
	cmd.Stdout = cconn
	cmd.Stderr = cconn
	sigs := make(chan error, 1)
	go monitorProcess(cmd, sigs)
	go func() {
		err := cmd.Run()
		sigs <- err
	}()
	err = <-sigs
	return
}

func monitorProcess(cmd *exec.Cmd, sigs chan error) {
	for {
		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			sigs <- nil
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
}
