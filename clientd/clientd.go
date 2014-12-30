package main

import (
	"flag"
	"github.com/akamel001/go-daemon"
	"github.com/akamel001/go-toml"
	"log"
	"os"
	"syscall"
	//"time"
)

var (
	signal = flag.String("s", "", `send signal to the daemon
    quit — graceful shutdown
    stop — fast shutdown
    reload — reloading the configuration file`)
)

func main() {
	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

	cntxt := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"SSG_CLIENT_DAEMON"},
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalln("Unable send signal to the daemon:", err)
		}
		daemon.SendCommands(d)
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
		log.Println("Error:", err)
	}
	log.Println("daemon terminated")
}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func worker() {
	for {
		config, err := toml.LoadFile("./clientd.conf")
		if err != nil {
			log.Println("Error ", err.Error())
		} else {

			configTree := config.Get("postgres").(*toml.TomlTree)
			user := configTree.Get("user").(string)
			password := configTree.Get("password").(string)
			port := configTree.Get("port").(int64)
			log.Println("User is ", user, ". Password is ", password, " port (", port, ")")
		}
		if _, ok := <-stop; ok {
			break
		}
	}
	//for {
	//log.Println("Sleeping for ", time.Second)
	// time.Sleep(time.Second)
	// if _, ok := <-stop; ok {
	//   log.Println("got stop signal")
	//   break
	// }
	//}

	// Jump back to done to exit
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
