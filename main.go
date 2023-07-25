package main

import (
	"fmt"

	"github.com/Arch-4ng3l/ChatForum/api"
	"github.com/Arch-4ng3l/ChatForum/storage"
)

func main() {
	psql := storage.NewPostgresDB()

	if psql == nil {
		fmt.Println("Couldn't Create Connection To Database")
	}

	api.NewAPIServer(":3333", psql).Run()
}
