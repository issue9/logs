// SPDX-License-Identifier: MIT

package main

import (
	"os"

	"github.com/issue9/logs/v2"
)

func main() {
	err := logs.InitFromXMLFile("./config.xml")
	if err != nil {
		//panic(err)
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	defer logs.Flush()

	logs.Info("INFO1")
	logs.Debugf("DEBUG %v", 1)
	logs.ERROR().Println("ERROR().Println")
}
