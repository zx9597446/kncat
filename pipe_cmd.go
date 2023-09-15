package main

import (
	"os/exec"
	"time"
)

func pipeCmd2Conn(strCmd string, cconn *CryptConn) (err error) {
	logf("pipe cmd: %s -> %s", strCmd, cconn.RemoteAddr().String())
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
	logf("pipe cmd done: %v", err)
	return
}

func monitorProcess(cmd *exec.Cmd, sigs chan error) {
	for {
		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			logf("process exited: %s", cmd.String())
			sigs <- nil
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
}
