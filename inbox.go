package main

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
	"log"
	"sort"
	"strconv"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	shell "github.com/stateless-minds/go-ipfs-api"
)

const (
	topicUpdateIssue = "update-issue"
)

type Inbox struct{
	app.Compo
	sh    *shell.Shell
	sub   *shell.PubSubSubscription
	issues []Issue
	alert *Alert
	veto bool
}

type Issue struct {
	ID          string      `mapstructure:"_id" json:"_id" validate:"uuid_rfc4122"`
	Title 		string 		`mapstructure:"title" json:"title" validate:"uuid_rfc4122"`
	Description string 		`mapstructure:"description" json:"description" validate:"uuid_rfc4122"`
	Solutions   []Solution  `mapstructure:"solutions" json:"solutions" validate:"uuid_rfc4122"`
	Veto 		bool		`mapstructure:"veto" json:"veto" validate:"uuid_rfc4122"`
}

type Solution struct {
	ID          string     `mapstructure:"_id" json:"_id" validate:"uuid_rfc4122"`
	Body 		string     `mapstructure:"body" json:"body" validate:"uuid_rfc4122"`
}

func newInbox() *Inbox {
	return &Inbox{}
}


func (c *Inbox) OnMount(ctx app.Context) {
	sh := shell.NewShell("localhost:5001")
	c.sh = sh

	c.subscribeToCreateIssueTopic(ctx)
	c.subscribeToUpdateIssueTopic(ctx)

	// c.DeleteIssues(ctx)
	c.FetchIssues(ctx)

	c.alert = newAlert()

	ctx.Handle("refresh", c.handleRefresh) // Registering action handler.
}

func (c *Inbox) handleRefresh(ctx app.Context, a app.Action) {
	ctx.Update()
}

func (c *Inbox) onChangeTab(ctx app.Context, e app.Event) {
	id := ctx.JSSrc().Get("id").String()
	app.Window().Get("document").Call("querySelector", ".email-item-unread").Get("classList").Call("remove", "email-item-unread")
	app.Window().Get("document").Call("querySelector", ".email-item-selected").Get("classList").Call("remove", "email-item-selected")

	app.Window().GetElementByID(id).Get("parentElement").Get("parentElement").Get("classList").Call("add", "email-item-unread")
	app.Window().GetElementByID(id).Get("parentElement").Get("parentElement").Get("classList").Call("add", "email-item-selected")

	app.Window().Get("document").Call("querySelector", ".inbox.active").Get("classList").Call("remove", "active")

	mainElem := app.Window().
    Get("document").
    Call("getElementById","main-"+id)

	mainElem.Get("classList").Call("add", "active")
}

func (c *Inbox) DeleteIssues(ctx app.Context) {
	ctx.Async(func() {
		err := c.sh.OrbitDocsDelete(dbNameIssue, "all")
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not delete issues")
			log.Fatal(err)
		}

		ctx.Dispatch(func(ctx app.Context) {
			c.issues = []Issue{}
		})

	})
}

func (c *Inbox) FetchIssues(ctx app.Context) {
	ctx.Async(func() {
		v, err := c.sh.OrbitDocsQuery(dbNameIssue, "all", "")
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not fetch issues")
			log.Fatal(err)
		}

		var vv []interface{}
		err = json.Unmarshal(v, &vv)
		if err != nil {
			log.Fatal(err)
		}

		var issues []Issue
		for _, ii := range vv {
			issue := Issue{}
			err = mapstructure.Decode(ii, &issue)
			if err != nil {
				log.Fatal(err)
			}

			issues = append(issues, issue)
			sort.SliceStable(issues, func(i, j int) bool {
				return issues[i].ID < issues[j].ID
			})

			if len(issue.Solutions) > 0 {
				sort.SliceStable(issue.Solutions, func(i, j int) bool {
					return issue.Solutions[i].ID < issue.Solutions[j].ID
				})
			}
		}

		ctx.Dispatch(func(ctx app.Context) {
			c.issues = issues
			log.Println(c.issues)
		})
	})
}

func (c *Inbox) subscribeToCreateIssueTopic(ctx app.Context) {
	ctx.Async(func() {
		topic := topicCreateIssue
		subscription, err := c.sh.PubSubSubscribe(topic)
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not subscribe to topic")
			log.Fatal(err)
		}
		c.sub = subscription
		c.subscriptionCreateIssue(ctx)
	})
}

func (c *Inbox) subscriptionCreateIssue(ctx app.Context) {
	ctx.Async(func() {
		defer c.sub.Cancel()
		// wait on pubsub
		res, err := c.sub.Next()
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not subscribe to topic")
			log.Fatal(err)
		}
		// Decode the string data.
		str := string(res.Data)
		log.Println("Subscriber of topic create-issue received message: " + str)
		ctx.Async(func() {
			c.subscribeToCreateIssueTopic(ctx)
		})

		i := Issue{}
		err = json.Unmarshal([]byte(str), &i)
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not subscribe to topic")
			log.Fatal(err)
		}

		ctx.Dispatch(func(ctx app.Context) {
			c.issues = append(c.issues, i)
		})
	})
}

func (c *Inbox) subscribeToUpdateIssueTopic(ctx app.Context) {
	ctx.Async(func() {
		topic := topicUpdateIssue
		subscription, err := c.sh.PubSubSubscribe(topic)
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not subscribe to topic")
			log.Fatal(err)
		}
		c.sub = subscription
		c.subscriptionUpdateIssue(ctx)
	})
}

func (c *Inbox) subscriptionUpdateIssue(ctx app.Context) {
	ctx.Async(func() {
		defer c.sub.Cancel()
		// wait on pubsub
		res, err := c.sub.Next()
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not subscribe to topic")
			log.Fatal(err)
		}
		// Decode the string data.
		str := string(res.Data)
		log.Println("Subscriber of topic update-issue received message: " + str)
		ctx.Async(func() {
			c.subscribeToUpdateIssueTopic(ctx)
		})

		i := Issue{}
		err = json.Unmarshal([]byte(str), &i)
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not subscribe to topic")
			log.Fatal(err)
		}

		ctx.Dispatch(func(ctx app.Context) {
			for _, iss := range c.issues {
				if iss.ID == i.ID {
					iss = i
				}
			}
		})
	})
}

func (c *Inbox) Render() app.UI {
	// log.Println(c.alert)
    return app.Div().Class("layout content pure-g").Body(
		app.If(c.alert != nil && c.alert.message != "", func() app.UI {
			return c.alert.Render()
		}),
		newNav(),
		app.Div().ID("list").Class("pure-u-1 inbox-list").Body(
			app.Range(c.issues).Slice(func(i int) app.UI {
				return app.If(i == 0, func() app.UI {
					return app.Div().ID("inbox-item-" + c.issues[i].ID).Class("email-item email-item-unread email-item-selected pure-g").Body(
						app.Div().Class("pure-u-1").Body(
							app.A().ID("inbox-" + c.issues[i].ID).Href("").Body(
								app.H4().Class("email-subject").Text(c.issues[i].Title),
							).OnClick(c.onChangeTab),
						),
					)
				}).Else(func() app.UI {
					return app.Div().ID("inbox-item-" + c.issues[i].ID).Class("email-item pure-g").Body(
						app.Div().Class("pure-u-1").Body(
							app.A().ID("inbox-" + c.issues[i].ID).Href("").Body(
								app.H4().Class("email-subject").Text(c.issues[i].Title),
							).OnClick(c.onChangeTab),
						),
					)
				})
			}),
    	),
		// Iterate over issues to display them...
		app.Range(c.issues).Slice(func(i int) app.UI {
			return app.If(i == 0, func() app.UI {
				return c.constructInboxItem(i, "active")	
			}).Else(func() app.UI {
				return c.constructInboxItem(i, "")
			})
		}),
	)
}

func (c *Inbox) constructInboxItem(index int, activeClass string) app.UI {
	return app.Div().ID("main-inbox-" + c.issues[index].ID).Class("pure-u-1 inbox " + activeClass).Body(
		app.Div().Class("email-content").Body(
			app.Div().Class("email-content-header pure-g").Body(
				app.Div().Class("pure-u-1").Body(
					app.H1().Class("email-content-title").Text(c.issues[index].Title),
				),
			),
			app.Div().Class("email-content-body").Body(
				app.P().Text(c.issues[index].Description),
			),
		),
		app.Div().Class("email-content").Body(
			app.Div().Class("email-content-header pure-g").Body(
				app.Div().Class("pure-u-2-3").Body(
					app.H2().Class("email-content-title").Text("Solutions"),
				),
				app.Div().Class("email-content-controls pure-u-1-3").Body(
					app.If(c.issues[index].Veto && len(c.issues[index].Solutions) > 0, func() app.UI {
						return app.Button().ID(strconv.Itoa(index)).Class("secondary-button pure-button").Text("Consensus").OnClick(c.onClickConsensus)
					}),
					app.If(!c.issues[index].Veto, func() app.UI {
						return app.Button().ID(strconv.Itoa(index)).Class("secondary-button pure-button").Text("Veto").OnClick(c.onClickVeto)
					}),
				),
			),
		),
		app.Range(c.issues[index].Solutions).Slice(func(n int) app.UI {
			return app.Div().Class("email-content").Body(
				app.Div().Class("email-content-body").Body(
					app.P().Text(c.issues[index].Solutions[n].Body),
				),
			)
		}),
		app.If(c.issues[index].Veto, func() app.UI {
			return c.constructSolutionForm(index)
		}),
	)
}

// func (c *Inbox) onClickVeto(ctx app.Context, e app.Event) {
// 	e.PreventDefault()
// 	c.veto = true
// }

// func (c *Inbox) onClickConsensus(ctx app.Context, e app.Event) {
// 	e.PreventDefault()
// 	c.veto = false
// }

func (c *Inbox) constructSolutionForm(id int) app.UI {
	return app.Div().Class("pure-u-1 form-box").Body(
		app.H1().Text("Suggest solution"),
		app.Form().ID(strconv.Itoa(id)).Class("pure-form").Body(
			app.FieldSet().Class("pure-group").Body(
				app.Textarea().ID("reply").Class("pure-input-1").Placeholder("How would you solve the issue?").Cols(50).Rows(10).Required(true),
			),
			app.Div().Class("email-content-controls pure-u-1").Body(
				app.Button().Class("pure-button pure-button-primary").Type("submit").Text("Reply"),
			),
		).OnSubmit(c.onSubmitSolution),
	)
}

func (c *Inbox) onSubmitSolution(ctx app.Context, e app.Event) {
	e.PreventDefault()

	index := ctx.JSSrc().Get("id").String()
	indexInt, err := strconv.Atoi(index)
	if err != nil {
		log.Fatal(err)
	}
	reply := app.Window().GetElementByID("reply").Get("value").String()

	id := uuid.New()

	s := Solution{
		ID:           id.String(),
		Body:         reply,
	}

	c.issues[indexInt].Solutions = append(c.issues[indexInt].Solutions, s)

	issue, err := json.Marshal(c.issues[indexInt])
	if err != nil {
		c.alert.send(ctx, AlertError, "Could not save solution")
		log.Fatal(err)
	}

	ctx.Async(func() {
		err = c.sh.OrbitDocsPut(dbNameIssue, issue)
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not save solution")
			log.Fatal(err)
		}

		err = c.sh.PubSubPublish(topicUpdateIssue, string(issue))
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not save solution")
			log.Fatal(err)
		}
	})
}

func (c *Inbox) onClickConsensus(ctx app.Context, e app.Event) {
	e.PreventDefault()

	index := ctx.JSSrc().Get("id").String()
	indexInt, err := strconv.Atoi(index)
	if err != nil {
		log.Fatal(err)
	}

	c.issues[indexInt].Veto = false

	issue, err := json.Marshal(c.issues[indexInt])
	if err != nil {
		c.alert.send(ctx, AlertError, "Could not save solution")
		log.Fatal(err)
	}

	ctx.Async(func() {
		err = c.sh.OrbitDocsPut(dbNameIssue, issue)
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not save solution")
			log.Fatal(err)
		}

		err = c.sh.PubSubPublish(topicUpdateIssue, string(issue))
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not save solution")
			log.Fatal(err)
		}
	})
}

func (c *Inbox) onClickVeto(ctx app.Context, e app.Event) {
	e.PreventDefault()

	index := ctx.JSSrc().Get("id").String()
	indexInt, err := strconv.Atoi(index)
	if err != nil {
		log.Fatal(err)
	}

	c.issues[indexInt].Veto = true
	
	issue, err := json.Marshal(c.issues[indexInt])
	if err != nil {
		c.alert.send(ctx, AlertError, "Could not save solution")
		log.Fatal(err)
	}

	ctx.Async(func() {
		err = c.sh.OrbitDocsPut(dbNameIssue, issue)
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not save solution")
			log.Fatal(err)
		}

		err = c.sh.PubSubPublish(topicUpdateIssue, string(issue))
		if err != nil {
			c.alert.send(ctx, AlertError, "Could not save solution")
			log.Fatal(err)
		}
	})
}
