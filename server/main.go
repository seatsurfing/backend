package main

import (
	"log"
	"os"
)

func main() {
	log.Println("Starting...")
	db := GetDatabase()
	a := GetApp()
	a.InitializeDatabases()
	a.InitializeRouter()
	a.InitializeAtlassianConnect()
	a.InitializeTimers()
	if GetConfig().PrintConfig {
		GetConfig().Print()
	}
	a.Run(GetConfig().PublicListenAddr)
	db.Close()
	os.Exit(0)
}
