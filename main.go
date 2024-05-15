package main

import (
	"log"
	"xxx-server/cmd"
)

var (
	version   string
	buildTime string
)

// @title       XXX SERVER API
// @version     0.0.1
// @description XXX系统后端 HTTP REST API
// @BasePath    /xxx-server/v1
func main() {
	log.Println("xxx-server start to run...", "version:", version, "buildTime:", buildTime)
	// log.Println("env:", os.Environ())
	cmd.Execute()
}
