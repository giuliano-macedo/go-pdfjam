package main

import (
	"os"

	"github.com/giuliano-macedo/go-pdfjam/internal/core/pdfjoin"
	"github.com/giuliano-macedo/go-pdfjam/internal/core/web"

	"github.com/skratchdot/open-golang/open"
	"github.com/sqweek/dialog"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		server, err := web.NewServer()
		if err != nil {
			panic(err)
		}

		err = server.Run(":8080")
		if err != nil {
			panic(err)
		}
	} else {
		file, err := dialog.File().Title("Save Result As").Filter("PDF Files", "pdf").Save()

		if err != nil {
			dialog.Message("error saving file: %s", err).Error()
			panic(err)
		}

		err = pdfjoin.Join(args, file)
		if err != nil {
			dialog.Message("error joining pdf files: %s", err).Error()
			panic(err)
		}
		err = open.Run(file)
		if err != nil {
			dialog.Message("error opening result file: %s", err).Error()
			panic(err)
		}
	}

}
