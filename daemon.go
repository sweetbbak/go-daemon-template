package main

import (
	"flag"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"syscall"
	"time"
)

var (
	_signal = flag.String("s", "", `Send signal to the daemon:
  quit — graceful shutdown
  stop — fast shutdown
  reload — reloading the configuration file`)
	execute = flag.String("e", "", "ask the daemon to execute a process")
	// SigSet  = flag.NewFlagSet("")
)

var (
	sock     = "unixsock"
	sockAddr = createSocket(sock)
)

func init() {
	flag.StringVar(_signal, "signal", "", "Send signal to the daemon")
	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(_signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(_signal, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(_signal, "reload"), syscall.SIGHUP, reloadHandler)

	// ./daemon stop VS ./daemon --signal stop
	// you could also add another FlagSet and decide to parse it here, to make use of sub-commands
	args := flag.Args()
	if len(args) >= 1 {
		switch args[0] {
		case "stop":
			*_signal = "stop"
		case "quit":
			*_signal = "quit"
		case "reload":
			*_signal = "reload"
		}
	}
}

func main() {
	cntxt := &daemon.Context{
		PidFileName: "sample.pid",
		PidFilePerm: 0644,
		LogFileName: "sample.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[go-daemon]"},
	}

	if *execute != "" {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatal(err)
		}

		// check if daemon is running, if so, we send a message to the server to be handled
		if d != nil {
			Client(sockAddr, *execute)
			os.Exit(0)
		}
	}

	// NOTE: if you dont cleanup pid file, this will give a nil pointer reference
	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalf("Unable send signal to the daemon: %s", err.Error())
		}
		if d != nil {
			daemon.SendCommands(d)
		} else {
			log.Println("Daemon is not currently running")
		}
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("daemon started")

	go worker()

	err = daemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	log.Println("daemon terminated")
}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func worker() {
LOOP:
	for {

		// serve and listen at our socket address
		Server(sockAddr)

		time.Sleep(time.Millisecond * 500) // this is work to be done by worker.
		select {
		case <-stop:
			break LOOP
		default:
		}
	}
	done <- struct{}{}
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	log.Println("configuration reloaded")
	return nil
}
