package main

import (
	"manager/config"
)

// @title FVRK API
// @version 0.1
// @description API
// @host 127.0.0.1:8000
// @BasePath /

var Config struct {
	Maria_ip         string
	Maria_user       string
	Maria_password   string
	Maria_database   string
	Mongo_ip         string
	Mongo_database   string
	Mongo_collection string
	Smtp_host        string
	Smtp_from        string
	Smtp_password    string
}

func main() {
	// "192.168.1.17:8002"
	// "tcp(192.168.1.17:8001)"
	// "172.18.0.9:27017"
	// "tcp(172.17.0.10:3306)"
	config.Init()
	HandleRequests()
}
