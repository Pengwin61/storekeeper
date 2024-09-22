package main

import (
	"os"
	"storekeeper/client"
	"storekeeper/database"
)

func main() {
	db := database.InitDB()
	tbt := os.Getenv("TELEGRAM_BOT_TOKEN")
	telega := client.NewClient(tbt)
	telega.Start(db)
}
