// SPDX-License-Identifier: MIT

package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/issue9/logs/v2"
	"github.com/issue9/logs/v2/config"
)

func main() {
	data, err := ioutil.ReadFile("./config.xml")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	cfg := &config.Config{}
	if err := xml.Unmarshal(data, cfg); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	l := logs.New()
	if err = l.Init(cfg); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	defer l.Flush()

	l.Info("INFO1")
	l.Debugf("DEBUG %v", 1)
	l.ERROR().Println("ERROR().Println")
}
