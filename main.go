package main

import "github.com/giuliano-macedo/go-pdfjam/internal/core/web"

func main() {
	server, err := web.NewServer()
	if err != nil {
		panic(err)
	}

	err = server.Run(":8080")
	if err != nil {
		panic(err)
	}

}
