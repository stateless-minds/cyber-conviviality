package main

import (
	"log"
	"net/http"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
    app.Route("/", func() app.Composer {
		return &Tutorials{}
	})

	app.RouteWithRegexp("^/tutorials/.*", func() app.Composer {
		return &Tutorials{}
	})

    app.Route("/compose", func() app.Composer {
		return &Compose{}
	})

	app.RouteWithRegexp("/inbox", func() app.Composer {
		return &Inbox{}
	})
    
	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:        "Cyber Conviviality",
		Description: "Learn decision-making by practice!",
		Styles: []string{
			"https://cdn.jsdelivr.net/npm/purecss@3.0.0/build/pure-min.css", // pureCSS
			"web/app.css",
		},
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
