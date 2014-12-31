package main

import (
	"flag"
	"github.com/akamel001/go-daemon"
	"github.com/akamel001/go-toml"
	"github.com/akamel001/ssg/libs"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
)

var (
	config_file = flag.String("c", "./clientd.conf", "What config to use for the client")
)

func main() {
	flag.Parse()

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
	var user, password, host string
	var port int64

	log.Printf("Loading config from %s", *config_file)
	config, err := toml.LoadFile(*config_file)

	if err != nil {
		log.Fatal("Error ", err.Error())
	} else {
		tree := config.Get("clientd").(*toml.TomlTree)
		user = tree.Get("user").(string)
		password = tree.Get("password").(string)
		port = tree.Get("port").(int64)
		host = tree.Get("host").(string)
		log.Println("User is ", user, ". Password is ", password, " port (", port, ")")
	}

	conn, err := net.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(int(port))))

	if err != nil {
		log.Println("Error ", err.Error())
	}

	for {
		msg := &ssg.DataPoint{
			Source:   proto.String("suck.a.dick.com"),
			Label:    proto.String("Impossibru"),
			IntValue: proto.Int64(12345),
		}
		line_msg, err := proto.Marshal(msg)
		if err != nil {
			log.Println("Failed to encode test message")
		}
		log.Printf("Sending message: %x", line_msg)
		conn.Write(line_msg)
		time.Sleep(1 * time.Second)
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
