package main

import (
	"log"
	"net/http"
	"os"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"uuid-gen/pages"
	"uuid-gen/utils"
)

const assetsDir = "web"

func main() {
	app.Route("/", func() app.Composer {
		return &pages.UUIDGen{}
	})

	app.RunWhenOnBrowser()

	// Define the handler with the required configuration
	handler := &app.Handler{
		Name:        "UUID Generator",
		Description: "Minimalist Online UUID Generator (v1, v4, v7)",
		Styles: []string{
			"https://cdn.jsdelivr.net/npm/halfmoon@1.1.1/css/halfmoon.min.css",
		},
		Title: "UUID Generator",
		Icon: app.Icon{
			Default: "/" + assetsDir + "/icon.svg",
		},
	}

	// Check if the "dist" argument is provided to generate the static site
	if len(os.Args) > 1 && os.Args[1] == "dist" {
		if err := app.GenerateStaticWebsite("dist", handler); err != nil {
			log.Fatal(err)
		}
		if err := utils.CopyAssetsToDist(assetsDir); err != nil {
			log.Fatal(err)
		}
		if err := utils.UpdateWasmExec(); err != nil {
			log.Fatal(err)
		}
		return
	}

	http.Handle("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
