package main

import (
	"flag"
	"github.com/akamel001/go-daemon"
	"github.com/akamel001/go-toml"
	"log"
	"os"
	"syscall"
	"net"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/akamel001/ssg/libs"
	//"io"
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
		Args:        []string{"SSG_SERVER_DAEMON"},
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
	//for {
		config, err := toml.LoadFile("./serverd.conf")
		if err != nil {
			log.Println("Error ", err.Error())
		} else {

			configTree := config.Get("postgres").(*toml.TomlTree)
			user := configTree.Get("user").(string)
			password := configTree.Get("password").(string)
			port := configTree.Get("port").(int64)
			log.Println("User is ", user, ". Password is ", password, " port (", port, ")")
			l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			log.Println("Listening for connections on port ", port)
			if err != nil {
				log.Fatal(err)
			}
			defer l.Close()
			for {
				// Wait for a connection.
				conn, err := l.Accept()
				if err != nil {
					log.Fatal(err)
				}

				go handle(conn)

				// go func(c net.Conn) {
				// 	log.Println("Received a new connection from ", c.RemoteAddr())
				// 	//log.Println("Received data ", c.data)
				// 	var buff = make([]byte, 30)
				// 	for {
				// 		readlen, ok := c.Read(buff)
				// 		if ok != nil {
				// 			log.Println("Error reading from socket ", ok)
				// 			break
				// 		}
				// 		if readlen == 0 {
				// 			log.Println("Connection closed by remote host")
				// 			break
				// 		}
				// 		//log.Println("Message: ", string(buff))
				// 		log.Printf("Got message: %x", buff)
				// 		buff = buff[:0]
				// 	}
			    
				// 	//fmt.Fscan(conn, &cmd)
    // 			//log.Println(fmt.Println("Message:", string(cmd)))
    // 			//log.Printf("%x", cmd)
				// 	// Echo all incoming data.
				// 	//io.Copy(c, c)
				// 	// Shut down the connection.
				// 	c.Close()
				// }(conn)

				//TODO: **** Need ti find a way to handle exit signal better **** 
				// if _, ok := <-stop; ok {
				// 	break
				// }		
			}
		}

//	}
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

func handle(c net.Conn){
	log.Println("Received a new connection from ", c.RemoteAddr())
	//TODO: *** create debug flag *** 
	//log.Println("Received data ", c.data) 

	//TODO: Message buffer is not exact and causes error when reading from socket 
	var buff = make([]byte, 1024)
	for {
		readlen, ok := c.Read(buff)
		if ok != nil {
			log.Println("Error reading from socket ", ok)
			break
		}
		if readlen == 0 {
			log.Println("Connection closed by remote host")
			break
		}
		//log.Println("Message: ", string(buff))
		msg := new(ssg.DataPoint)

		err := proto.Unmarshal(buff, msg)
		if err != nil{
			log.Println("Error received while trying to unmarshal message ", err)
		}
		log.Printf("Got message: %v", msg)
	}
  
	//fmt.Fscan(conn, &cmd)
	//log.Println(fmt.Println("Message:", string(cmd)))
	//log.Printf("%x", cmd)
	// Echo all incoming data.
	//io.Copy(c, c)
	// Shut down the connection.

	//Need a better way to handle closing the connection 
	//c.Close()
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
