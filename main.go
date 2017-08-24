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
	"fmt"
	"log"
	"os"
	"runtime"
)

type options struct {
	BindAddr    string
	LogFile     string
	ScoreDBFile string
}

func main() {
	conf_file_path := getConfFilePath()
	conf := readConfiguration(conf_file_path)

	serve(conf)
}

/* TODO account for all available Go platforms */
func getConfFilePath() string {
	if runtime.GOOS == "windows" {
		/* TODO figure out where Windows keeps its config files */
		fmt.Fprintf(os.Stderr, "No Windows support at this time\n")
		os.Exit(1)
	} else if runtime.GOOS == "freebsd" {
		return "/usr/local/etc/linux-scoring-engine.json"
	} else if runtime.GOOS == "darwin" {
		/* TODO make sure this is accurate */
		return "/Library/linux-scoring-engine.json"
	}

	return "/etc/linux-scoring-engine.json"
}

func bailIfFail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
