package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func serve(opts options) {
	ln, err := net.Listen("tcp", opts.BindAddr)
	bailIfFail(err)

	log := initLog(opts.LogFile, opts.BindAddr)

	for {
		conn, _ := ln.Accept()

		go handleConnection(conn, log)
	}
}

func initLog(logpath string, addr string) *os.File {
	file, err := os.OpenFile(logpath, os.O_RDWR|os.O_APPEND, 0600)
	bailIfFail(err)
	_, err = fmt.Fprintf(file, "Initializing cyberpatriot scoring server: bound to \"%s\"\n", addr)
	bailIfFail(err)

	return file
}

/*
 * proper request is in JSON, like below:
 *
 * {
 * 	"service":"sshd",
 * 	"pointchange":2,
 * 	"reason":"disabled auth without keys"
 * }
 */
type parsed_request struct {
	Service     string
	PointChange int
	Reason      string
}

func handleConnection(conn net.Conn, log *os.File) {
	defer conn.Close()

	fmt.Fprintf(log, "Incoming connection from %s: ", conn.RemoteAddr())

	readbuffer, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Fprintln(conn, "server read() failure: resend request")
		fmt.Fprintln(log, "server read() failure: resend request")
		return
	}

	_, err = parseRequest(readbuffer)
	if err != nil {
		fmt.Fprintln(conn, err)
	}
}

func parseRequest(request_buffer []byte) (parsed_request, error) {
	parsed := parsed_request{}
	err := json.Unmarshal(request_buffer, &parsed)
	if err != nil {
		return parsed, err
	}

	return parsed, nil
}
