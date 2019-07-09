package main

import (
	"log"
	"math/rand"
	"time"
)

func init() {
	log.Println("Initiating...")

	// Random seed
	rand.Seed(time.Now().UnixNano())

	// Setup data directories
	setupDataDirs()
}

func main() {
	// Connecting to datasources
	dbmgr = getDBmgr()
	defer dbmgr.Close()

	pubsub = getPubsub()
	defer pubsub.Close()

	// Migrate
	dbmgr.Migrate()

	log.Println("Starting Server...")
	srv := CreateServer()
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
