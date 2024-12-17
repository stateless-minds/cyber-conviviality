package main

import (
    "strings"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Nav struct{
	app.Compo
}


func newNav() *Nav {
	return &Nav{}
}

func (c *Nav) OnMount(ctx app.Context) {
    path := app.Window().URL().Path
    path = strings.ReplaceAll(path, "/", "")
    
    if path == "" {
        path = "tutorials"
        
    }

    app.Window().GetElementByID(path).Get("parentElement").Get("classList").Call("add", "pure-menu-selected")

}

func (n *Nav) Render() app.UI {
    return app.Div().ID("nav").Class("pure-u").Body(
        app.A().Href("/").ID("menuLink").Class("nav-menu-button").Text("Cyber Conviviality"),
        app.Div().Class("nav-inner").Body(
            app.Div().Class("pure-menu").Body(
                app.Ul().Class("pure-menu-list").Body(
                    app.Li().Class("pure-menu-item").Body(
                        app.A().ID("compose").Href("/compose").Class("primary-button pure-button").Text("Compose"),
                    ),
                    app.Li().Class("pure-menu-item").Body(
                        app.A().ID("tutorials").Href("/tutorials/#1").Class("pure-menu-link").Text("Tutorials"),
                    ).Style("text-align", "left"),
                    app.Li().Class("pure-menu-item").Body(
                        app.A().ID("inbox").Href("/inbox").Class("pure-menu-link").Text("Inbox"),
                    ).Style("text-align", "left"),
                ),
            ),
        ),
    )
}
