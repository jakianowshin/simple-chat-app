package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
)

var serverHost = flag.String("server_host", "localhost", "host address for server to listen")
var serverPort = flag.String("port", "9091", "port on which server listen")

func StartServerApp(ip, port string) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"ip": ip, "port": port}).Error("cannot_start_the_server")
		panic("cannot start the server, err: " + err.Error())
	}
	log.WithFields(log.Fields{"host": ip, "port": port}).Info("server_started")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.WithError(err).Error("error_while_accepting_connection")
		}
		log.WithField("remote_add", conn.RemoteAddr()).Info("new_connection")
		go startMirrorChat(conn)
	}
}

func startMirrorChat(conn net.Conn) {
	if conn == nil {
		return
	}
	defer conn.Close()
	tmp := make([]byte, 512)
	buf := make([]byte, 0, 4096)
	remoteAddr := conn.RemoteAddr()
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err == io.EOF {
				log.WithField("remote_addr", remoteAddr).Info("eof_for_connection")
				fmt.Printf("[%s] >>>: [last_message] %s\n", remoteAddr, tmp)
				buf = append(buf, []byte("[last_message]")...)
				buf = append(buf, tmp[:n]...)

				break
			}
			log.WithError(err).WithField("remote_addr", conn.RemoteAddr()).Error("error_while_reading")
			return
		}
		buf = append(buf, []byte("|")...)
		buf = append(buf, tmp[:n]...)
		if err := writeToConnection(conn, buf); err != nil {
			log.WithError(err).WithField("remote_addr", remoteAddr).Error("write_error")
			return
		}
		fmt.Printf("[%s] >>>: %s\n", remoteAddr, string(tmp[:n]))
	}
}

func writeToConnection(conn net.Conn, buf []byte) error {
	if conn == nil {
		return fmt.Errorf("nil_connection: cannot_write")
	}
	if len(buf) == 0 {
		return nil
	}
	total := 0
	for {
		n, err := conn.Write(buf[total:])
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			return nil
		}
		if total+n >= len(buf) {
			return nil
		}
		total += n
	}
	return nil
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)

}
func main() {
	flag.Parse()
	StartServerApp(*serverHost, *serverPort)
}
