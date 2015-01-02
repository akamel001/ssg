package main

import (
	"flag"
	"os"
	"kylelemons.net/go/daemon"
	"github.com/akamel001/go-toml"
	"fmt"
	"net"
	//"time"
)

var (
	loglvl = daemon.LogLevelFlag("loglvl")
	//log = daemon.LogFileFlag("log", 7777777)
	//delay = flag.Duration("delay", time.Nanosecond, "Delay between restarts")
	fork  = daemon.ForkPIDFlags("fork", "pidfile", "serverd.pid")
	config_file = flag.String("config", "./serverd.conf", "What config to use for the server")
	config, err = toml.LoadFile(*config_file)
	port_obj  = daemon.ListenFlag("port", "tcp", fmt.Sprintf(":%d", config.Get("postgres.port").(int64)), "port")
)

func main() {
	flag.Parse()
	daemon.Verbose.Printf("Command-line: %q", os.Args)
	//daemon.RedirectStdout = false
	//daemon.Verbose.Printf("Logging to file %s", *log.String())
	fork.Fork()

	if err != nil {
		daemon.Fatal.Printf("Failed to load config: ", err.Error())
	} else {
		daemon.Verbose.Printf("Loaded config file %s", *config_file)

		port, err := port_obj.Listen()
		if err != nil {
			daemon.Fatal.Printf("listen: %s", err)
		}

		go func() {
			for {
				conn, err := port.Accept()
				if err == daemon.ErrStopped {
					break
				}
				if err != nil {
					daemon.Error.Printf("accept: %s", err)
				}
				go handle(conn)
			}
			daemon.Verbose.Printf("Serve loop exited")
		}()
		// go func() {
		// 	time.Sleep(*delay)
		// 	*delay = 24 * time.Hour
		// 	daemon.Restart(15 * time.Second)
		// }()
		daemon.Run()
	}
}

func handle(c net.Conn){
	//daemon.Info.Printf("Received a new connection from ", c.RemoteAddr())
	defer c.Close()
	//TODO: *** create debug flag *** 
	//log.Println("Received data ", c.data) 

	//TODO: Message buffer is not exact and causes error when reading from socket 
	var buff = make([]byte, 1024)
	for {
		readlen, ok := c.Read(buff)
		if ok != nil {
			daemon.Warning.Printf("Error reading from socket %s", ok)
			break
		}
		if readlen == 0 {
			daemon.Warning.Printf("Connection closed by remote host")
			break
		}
		//log.Println("Message: ", string(buff))
		//msg := new(ssg.DataPoint)
		//err := proto.Unmarshal(buff, msg)
		// if err != nil{
		// 	daemon.Info.Printf("Error received while trying to unmarshal message ", err)
		// }
		c.Write(buff)
		daemon.Verbose.Printf("Got message: %v", buff)
	}
	daemon.Verbose.Printf("Closed handle!")
}  

