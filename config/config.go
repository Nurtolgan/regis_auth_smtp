package config

import (
	"encoding/json"
	"manager/debugger"
	"os"
)

var Config struct {
	Maria_ip         string
	Maria_user       string
	Maria_password   string
	Maria_database   string
	Mongo_ip         string
	Mongo_database   string
	Mongo_collection string
	Smtp_host        string
	Smtp_port        int
	Smtp_from        string
	Smtp_password    string
}

func Init() {
	file, err := os.ReadFile("/home/config.json")
	debugger.CheckError("Read File", err)
	debugger.CheckError("Unmarshal", json.Unmarshal(file, &Config))
}
