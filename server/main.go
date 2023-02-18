package main

import (
	"log"
	"os"
	"path/filepath"
)

var _productVersion = ""

func GetProductVersion() string {
	if _productVersion == "" {
		path, _ := filepath.Abs("./res/version.txt")
		data, err := os.ReadFile(path)
		if err != nil {
			return "UNKNOWN"
		}
		_productVersion = string(data)
	}
	return _productVersion
}

func main() {
	log.Println("Starting...")
	log.Println("Seatsurfing Backend Version " + GetProductVersion())
	db := GetDatabase()
	a := GetApp()
	a.InitializeDatabases()
	a.InitializeDefaultOrg()
	a.InitializeRouter()
	a.InitializeTimers()
	if GetConfig().PrintConfig {
		GetConfig().Print()
	}
	a.Run(GetConfig().PublicListenAddr)
	db.Close()
	os.Exit(0)
}
