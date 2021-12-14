package main

import (
	"bufio"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
)

var serverHost = flag.String("server_host", "localhost", "server host to connect")
var serverPort = flag.String("serverPort", "9091", "server serverPort to connect")

func startClient(host, port string) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.WithError(err).Error("cannot_connect_to_server")
	}
	go readFromServer(conn)
	reader := bufio.NewReader(os.Stdin)
	for {
		buff, err := reader.ReadString('\n')

		err = writeToConnection(conn, []byte(buff))
		if err != nil {
			log.WithError(err).Error("write_failed")
			return
		}
	}
}

func readFromServer(conn io.Reader) {
	buff := make([]byte, 2048)

	for {
		n, err := conn.Read(buff)
		if err != nil && err != io.EOF {
			log.WithError(err).Error("error_while_reading_from_server")
			return
		}
		fmt.Printf(">>>: %s\n", string(buff[:n]))
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

func main() {
	flag.Parse()
	startClient(*serverHost, *serverPort)
}
