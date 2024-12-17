package main

import (
	"time"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	AlertError = "error"
	AlertSuccess = "success"
)

type Alert struct {
	app.Compo
	message string
	level string
}

func newAlert() *Alert {
	return &Alert{}
}

func (c *Alert) send(ctx app.Context, level string, msg string) {
	c.message = level + ": "+ msg
	c.level = level

	ctx.Async(func() {
		time.Sleep(5 * time.Second)
		c.message = ""
		c.level = ""
		ctx.NewAction("refresh")
	})
}

func (c *Alert) Render() app.UI {
	return app.Div().Class("container alert "+c.level).Body(
		app.Text(c.message),
	)
}