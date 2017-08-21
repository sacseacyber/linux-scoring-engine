/*
 * Copyright (c) 2017, Mark Rogers
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright notice,
 *   this list of conditions and the following disclaimer in the documentation
 *   and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
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
	file, err := os.OpenFile(logpath, os.O_RDWR|os.O_APPEND|os.O_CREATE,
		0644)
	bailIfFail(err)

	_, err = fmt.Fprintf(file, "Initializing cyberpatriot scoring server: bound to \"%s\"\n", addr)
	bailIfFail(err)

	return file
}

/*
 * proper request is in JSON, like below:
 *
 * {
 *	"reqtype":"put",
 * 	"service":"sshd",
 * 	"pointchange":2,
 * 	"reason":"disabled auth without keys"
 * }
 */
type parsed_request struct {
	Reqtype     string
	Service     string
	PointChange int
	Reason      string
}

func handleConnection(conn net.Conn, log *os.File) {
	defer conn.Close()

	fmt.Fprintf(log, "Incoming connection from %s: ", conn.RemoteAddr())

	sockread_buffer, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Fprintln(conn, "failure: server error: resend request")
		fmt.Fprintln(log, "failure: server error: resend request")
		return
	}

	reqdata, err := extractRequestData(sockread_buffer)
	if err != nil {
		fmt.Fprintln(conn, "failure: client error:", err)
		fmt.Fprintln(log, "failure: client error:", err)
		return
	}

	_, err = getRequestType(reqdata.Reqtype)
	if err != nil {
		fmt.Fprintln(conn, "failure: invalid request type")
		fmt.Fprintln(log, "failure: invalid request type")
		return
	}

	fmt.Fprintln(conn, "success")
	fmt.Fprintln(log, "success")
}

func extractRequestData(request_buffer []byte) (parsed_request, error) {
	parsed := parsed_request{}
	err := json.Unmarshal(request_buffer, &parsed)

	return parsed, err
}

func getRequestType(reqtype string) (string, error) {
	if strings.ToUpper(reqtype) == "GET" || strings.ToUpper(reqtype) == "PUT" {
		return reqtype, nil
	}

	return "", nil
}
