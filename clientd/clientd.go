package main

import (
	"flag"
	"kylelemons.net/go/daemon"
	"github.com/akamel001/go-toml"
	"github.com/akamel001/ssg/libs"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"strconv"
	"time"
)

func init(){
	daemon.RedirectStdout = false
}

var (
	loglvl = daemon.LogLevelFlag("loglvl")
	log = daemon.LogFileFlag("log", 0644)
	fork  = daemon.ForkPIDFlags("fork", "pidfile", "clientd.pid")
	config_file =  flag.String("config", "./clientd.conf", "What config to use for the client")
)

func main() {
	flag.Parse()

	daemon.Verbose.Printf("Command-line: %q", os.Args)

	if len(os.Args) < 2 {
		daemon.Info.Printf("Daemon started without any arguments. Running in foreground.")
	}
	daemon.Info.Printf("config file: %s", *config_file)
	fork.Fork()

	go send_metrics(*config_file)
	daemon.Run()
}

func send_metrics(config_file string) {
	var user, password, host string
	var port int64
	config, err := toml.LoadFile(config_file)

	daemon.Info.Printf("Loading config from %s", config_file)

	if err != nil {
		daemon.Fatal.Printf("Error ", err.Error())
	} else {
		tree := config.Get("clientd").(*toml.TomlTree)
		user = tree.Get("user").(string)
		password = tree.Get("password").(string)
		port = tree.Get("port").(int64)
		host = tree.Get("host").(string)
		daemon.Info.Printf("User: %s \tPassword: %s\tHost: %s\tPort: %d", user, password, host, port)
	}

	conn, err := net.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(int(port))))

	if err != nil {
		daemon.Fatal.Printf("Error while dialing connection %s", err.Error())
	}

	for {
		msg := &ssg.DataPoint{
			Source:   proto.String("suck.a.dick.com"),
			Label:    proto.String("Impossibru"),
			IntValue: proto.Int64(12345),
		}
		line_msg, err := proto.Marshal(msg)
		if err != nil {
			daemon.Warning.Printf("Failed to encode test message")
		}
		daemon.Info.Printf("Sending message: %x", line_msg)
		conn.Write(line_msg)
		time.Sleep(1 * time.Second)
	}
	daemon.Verbose.Printf("Serve loop exited")
}

