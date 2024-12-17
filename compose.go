package main

import (
    "encoding/json"
    "log"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/google/uuid"
	shell "github.com/stateless-minds/go-ipfs-api"
)

const dbNameIssue = "issue"

const (
	topicCreateIssue = "create-issue"
)

type Compose struct{
	app.Compo
	sh    *shell.Shell
	alert *Alert
}

func (c *Compose) OnMount(ctx app.Context) {
	sh := shell.NewShell("localhost:5001")
	c.sh = sh
	c.alert = newAlert()

	ctx.Handle("refresh", c.handleRefresh) // Registering action handler.
}

func (c *Compose) handleRefresh(ctx app.Context, a app.Action) {
	ctx.Navigate("inbox")
}

func (c *Compose) Render() app.UI {
    return app.Div().Class("layout content pure-g").Body(
		app.If(c.alert != nil && c.alert.message != "", func() app.UI {
			return c.alert.Render()
		}),
		app.Div().Class("pure-u-1 form-box compose").Body(
			newNav(),
			app.H1().Text("Post issue"),
			app.Form().Class("pure-form").Body(
				app.FieldSet().Class("pure-group").Body(
					app.Input().ID("title").Type("text").Class("pure-input-1-2").Placeholder("Issue Title").Required(true),
					app.Textarea().ID("desc").Class("pure-input-1-2").Placeholder("Issue Description").Rows(10).Required(true),
				),
				app.Button().Class("pure-button pure-input-1-2 pure-button-primary").Type("submit").Text("Submit"),
			).OnSubmit(c.onSubmitIssue),
    	),
	)
}

func (c *Compose) onSubmitIssue(ctx app.Context, e app.Event) {
	e.PreventDefault()
	title := app.Window().GetElementByID("title").Get("value").String()
	desc := app.Window().GetElementByID("desc").Get("value").String()

	id := uuid.New()

	i := Issue{
		ID:           id.String(),
		Title:        title,
		Description:  desc,
		Veto: 		  true,
	}

	issue, err := json.Marshal(i)
	if err != nil {
		c.alert.send(ctx, AlertError, "Could not save issue")
		log.Fatal(err)
	}

	ctx.Async(func() {
		err = c.sh.OrbitDocsPut(dbNameIssue, issue)
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not save issue")
			log.Fatal(err)
		}
		err = c.sh.PubSubPublish(topicCreateIssue, string(issue))
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not save issue")
			log.Fatal(err)
		}

		ctx.Dispatch(func(ctx app.Context) {
			c.alert.send(ctx, AlertSuccess, "Issue submitted")
		})
	})
}
