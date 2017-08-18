package main

import (
	"fmt"
	"os"
)

type options struct {
	BindAddr string
	LogFile  string
}

func main() {
	conf := readConfiguration("/etc/scored.json")

	serve(conf)
}

func bailIfFail(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
