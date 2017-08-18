package main

import (
	"encoding/json"
	"os"
)

func readConfiguration(path string) options {
	configfile, err := os.Open(path)
	bailIfFail(err)

	decoder := json.NewDecoder(configfile)
	opts := options{}
	err = decoder.Decode(&opts)

	bailIfFail(err)
	return opts
}
