package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

var (
	// application name
	Name = ""
	// application version string
	Version = ""
	// commit
	Commit = ""
	// build tags
	BuildTags = ""
	// application variable
	udpNetwork = "udp"
	tcpNetwork = "tcp"
	udpBufSize = 64 * 1024
)

var (
	logger = log.New(os.Stderr, "", 0)
)

func logf(f string, v ...interface{}) {
	if config.Verbose {
		logger.Output(2, fmt.Sprintf(f, v...))
	}
}

var config struct {
	Help         bool
	Verbose      bool
	Listen       bool
	Port         int
	Network      string
	Web          bool
	Command      bool
	Host         string
	Secret       string
	CryptoMethod string
}

func init() {
	flag.IntVar(&config.Port, "p", 4000, "host port to connect or listen")
	flag.BoolVar(&config.Help, "help", false, "print this help")
	flag.BoolVar(&config.Verbose, "v", true, "verbose mode")
	flag.BoolVar(&config.Listen, "l", false, "listen mode")
	flag.BoolVar(&config.Command, "e", false, "shell mode")
	flag.StringVar(&config.Network, "n", "tcp", "network protocol")
	flag.StringVar(&config.Host, "h", "0.0.0.0", "host addr to connect or listen")
	flag.StringVar(&config.Secret, "k", "", "secret key to crypt")
	flag.StringVar(&config.CryptoMethod, "m", "rc4", "crypto method: rc4 aes256cfb")
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Println(`
usage: kncat [-l] [-v] [-p port] [-n tcp] -k secret
options:`)
	flag.PrintDefaults()
}

func listen(network, host string, port int, command bool) {
	listenAddr := net.JoinHostPort(host, strconv.Itoa(port))
	listener, err := net.Listen(network, listenAddr)
	logf("Listening on: %s://%s", network, listenAddr)
	if err != nil {
		logf("Listen failed: %s", err)
		return
	}
	conn, err := listener.Accept()
	if err != nil {
		logf("Accept failed: %s", err)
		return
	}
	logf("Connection received: %s", conn.RemoteAddr())
	if command {
		var shell string
		switch runtime.GOOS {
		case "linux":
			shell = "/bin/sh"
		case "freebsd":
			shell = "/bin/csh"
		case "windows":
			shell = "cmd.exe"
		default:
			shell = "/bin/sh"
		}
		cmd := exec.Command(shell)
		cconn := NewCryptConn(conn, config.CryptoMethod, []byte(config.Secret))
		cmd.Stdin = cconn
		cmd.Stdout = cconn
		cmd.Stderr = cconn
		cmd.Run()
		defer conn.Close()
		logf("Closed: %s", conn.RemoteAddr())
	} else {
		go func(c net.Conn) {
			io.Copy(os.Stdout, c)
			c.Close()
			logf("Closed: %s", conn.RemoteAddr())
			os.Exit(0)
		}(conn)
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
			io.Copy(conn, os.Stdin)
		}
	}
}

func listenPacket(network, host string, port int, command bool) {
	listenAddr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.ListenPacket(network, listenAddr)
	if err != nil {
		logf("Listen failed: %s", err)
		return
	}
	logf("Listening on: %s://%s", network, listenAddr)
	defer func(c net.PacketConn) {
		logf("\nClosed udp listen")
		c.Close()
		os.Exit(0)
	}(conn)
	buf := make([]byte, udpBufSize)
	n, addr, err := conn.ReadFrom(buf)
	if n == 0 || err == io.EOF {
		return
	}
	logf("Connection received : %s", addr.String())
	fmt.Fprint(os.Stdout, string(buf))
}

func dial(network, host string, port int, command bool) {
	dailAddr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.Dial(network, dailAddr)
	if err != nil {
		logf("Dail failed: %s", err)
		return
	}
	logf("Dialed host: %s://%s", network, dailAddr)
	defer func(c net.Conn) {
		logf("Closed: %s", dailAddr)
		c.Close()
	}(conn)
	if command {
		var shell string
		switch runtime.GOOS {
		case "linux":
			shell = "/bin/sh"
		case "freebsd":
			shell = "/bin/csh"
		case "windows":
			shell = "cmd.exe"
		default:
			shell = "/bin/sh"
		}
		cmd := exec.Command(shell)
		cconn := NewCryptConn(conn, config.CryptoMethod, []byte(config.Secret))
		cmd.Stdin = cconn
		cmd.Stdout = cconn
		cmd.Stderr = cconn
		cmd.Run()
	} else {
		go io.Copy(os.Stdout, conn)
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
	}
}

func main() {
	if config.Help {
		flag.Usage()
		return
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigs
		logf("Exited")
		os.Exit(0)
	}()

	// Listen
	if config.Listen {
		switch config.Network {
		case udpNetwork:
			listenPacket(config.Network, config.Host, config.Port, config.Command)
		case tcpNetwork:
			listen(config.Network, config.Host, config.Port, config.Command)
		default:
			panic("no target network protocol")
		}
		// Dial
	} else {
		dial(config.Network, config.Host, config.Port, config.Command)
	}
}
