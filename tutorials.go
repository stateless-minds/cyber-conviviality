package main

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Tutorials struct{
	app.Compo
}

func (c *Tutorials) OnMount(ctx app.Context) {
	id := app.Window().URL().Fragment
	if id == "" {
		id = "1"
		ctx.Navigate("tutorials/#1")
	}

	app.Window().GetElementByID("tutorial-" + id).Get("parentElement").Get("parentElement").Get("classList").Call("add", "email-item-unread")
	app.Window().GetElementByID("tutorial-" + id).Get("parentElement").Get("parentElement").Get("classList").Call("add", "email-item-selected")
}

func (c *Tutorials) OnNav(ctx app.Context) {
    id := app.Window().URL().Fragment
	if id == "" {
		id = "1"
		ctx.Navigate("tutorials/#1")
	}

	c.changeTab(id)
}

func (c *Tutorials) changeTab(id string) {
	app.Window().Get("document").Call("querySelector", ".email-item-unread").Get("classList").Call("remove", "email-item-unread")
	app.Window().Get("document").Call("querySelector", ".email-item-selected").Get("classList").Call("remove", "email-item-selected")

	app.Window().GetElementByID("tutorial-" + id).Get("parentElement").Get("parentElement").Get("classList").Call("add", "email-item-unread")
	app.Window().GetElementByID("tutorial-" + id).Get("parentElement").Get("parentElement").Get("classList").Call("add", "email-item-selected")

	app.Window().Get("document").Call("querySelector", ".tutorial.active").Get("classList").Call("remove", "active")

	mainElem := app.Window().
    Get("document").
    Call("getElementById","main-tutorial-"+id)

	mainElem.Get("classList").Call("add", "active")
}

func (c *Tutorials) Render() app.UI {
    return app.Div().Class("layout content pure-g").Body(
		newNav(),
		app.Div().ID("list").Class("pure-u-1 tutorial-list").Body(
			app.Div().Class("email-item pure-g").Body(
				app.Div().Class("pure-u-1").Body(
					app.A().ID("tutorial-1").Href("#1").Body(
						app.H4().Class("email-subject").Text("Learn decision-making by practice"),
					),
				),
			),
			app.Div().Class("email-item pure-g").Body(
				app.Div().Class("pure-u-1").Body(
					app.A().ID("tutorial-2").Href("#2").Body(
						app.H4().Class("email-subject").Text("In memory of Ivan Illich"),
					),
				),
			),
			app.Div().Class("email-item pure-g").Body(
				app.Div().Class("pure-u-1").Body(
					app.A().ID("tutorial-3").Href("#3").Body(
						app.H4().Class("email-subject").Text("Redefining progress"),
					),
				),
			),
			app.Div().Class("email-item pure-g").Body(
				app.Div().Class("pure-u-1").Body(
					app.A().ID("tutorial-4").Href("#4").Body(
						app.H4().Class("email-subject").Text("The scenario"),
					),
				),
			),
			app.Div().Class("email-item pure-g").Body(
				app.Div().Class("pure-u-1").Body(
					app.A().ID("tutorial-5").Href("#5").Body(
						app.H4().Class("email-subject").Text("Features"),
					),
				),
			),
			app.Div().Class("email-item pure-g").Body(
				app.Div().Class("pure-u-1").Body(
					app.A().ID("tutorial-6").Href("#6").Body(
						app.H4().Class("email-subject").Text("Details"),
					),
				),
			),    
		),
		app.Div().ID("main-tutorial-1").Class("pure-u-1 tutorial active").Body(
			app.Div().Class("email-content").Body(
				app.Div().Class("email-content-header pure-g").Body(
					app.Div().Class("pure-u-1").Body(
						app.H1().Class("email-content-title").Text("Learn decision-making by practice"),
					),
				),
				app.Div().Class("email-content-body").Body(
					app.P().Text("Representative democracy and consumer culture deprived us from the learning process of decision making."),
					app.P().Text("All instutions are hierarchical and train us to listen, behave and comply."),
					app.P().Text("We have no say in decisions on any level from school, to work, to local council, to national and world scale."),
					app.P().Text("We are deprived of our very nature to participate and have a say in our common daily matters."),
					app.P().Text("Spaces which are by definition and history a commons such as our cities are de facto a state ownership."),
					app.P().Text("Decision making is the same as any other ability, you need to practice to get good at it."),
					app.P().Text("The purpose of the decision making simulator is to foster constructive collaboration, tolerance and empathy."),
					app.P().Text("Decision making is an art. Learn to participate, collaborate, swallow up your ego, compromise and reach consensus for the greater good."),
				),
			),
		),
		app.Div().ID("main-tutorial-2").Class("pure-u-1 tutorial").Body(
			app.Div().Class("email-content").Body(
				app.Div().Class("email-content-header pure-g").Body(
					app.Div().Class("pure-u-1").Body(
						app.H1().Class("email-content-title").Text("In memory of Ivan Illich"),
					),
				),
				app.Div().Class("email-content-body").Body(
					app.P().Text("Do you know what are convivial tools?"),
					app.P().Text("Have you heard of Ivan Illich?"),
					app.A().Href("https://monoskop.org/Ivan_Illich").Text("Ivan Illich"),
					app.P().Text("Once you are familiar with the concept you will be able to apply it successfully in finding convivial solutions to everyday problems."),
				),
			),
		),
		app.Div().ID("main-tutorial-3").Class("pure-u-1 tutorial").Body(
			app.Div().Class("email-content").Body(
				app.Div().Class("email-content-header pure-g").Body(
					app.Div().Class("pure-u-1").Body(
						app.H1().Class("email-content-title").Text("Redefining progress"),
					),
				),
				app.Div().Class("email-content-body").Body(
					app.P().Text("We are taught that progress is always something bigger, better and more powerful."),
					app.P().Text("This was the result of the growth paradigm. But the side effects are more visible than ever."),
					app.P().Text("From planned obsolescence to e-waste, to centralized unaccountable systems we are more remote than ever from the tools that are supposed to serve us."),
					app.P().Text("The production of electric cars and the rise of crypto currencies leads to the need for hundreds more nuclear power plants."),
					app.P().Text("The high-tech mania goes hand in hand with high power consumption and vice-versa."),
					app.P().Text("But we have low-tech convivial tools that do just what they are meant to and have no side effects. Things like bicycles, p2p apps, mobile homes, permaculture, autonomous cities etc."),
					app.P().Text("Without falling for gigantism, newness and value measuring we can have a completely new and unmeasarable definition of progress such as human satisfaction, autonomy, clean air, free time, spare space."),
				),
			),
		),
		app.Div().ID("main-tutorial-4").Class("pure-u-1 tutorial").Body(
			app.Div().Class("email-content").Body(
				app.Div().Class("email-content-header pure-g").Body(
					app.Div().Class("pure-u-1").Body(
						app.H1().Class("email-content-title").Text("The scenario"),
					),
				),
				app.Div().Class("email-content-body").Body(
					app.P().Text("Imagine for a moment that we have no countries, institutions, politicians, laws, money or property. All Earth's resources are commons."),
					app.P().Text("What would you do in a world starting from blank paper with the current level of technology and knowledge?"),
					app.P().Text("Would you repeat the same mistakes civilization did or would you take a different path in every aspect?"),
					app.P().Text("Let's find out and learn along the way!"),
				),
			),
		),
		app.Div().ID("main-tutorial-5").Class("pure-u-1 tutorial").Body(
			app.Div().Class("email-content").Body(
				app.Div().Class("email-content-header pure-g").Body(
					app.Div().Class("pure-u-1").Body(
						app.H1().Class("email-content-title").Text("Features"),
					),
				),
				app.Div().Class("email-content-body").Body(
					app.P().Text("- Post real world issues"),
					app.P().Text("- Post suggestions how to solve them"),
					app.P().Text("- Practice negotiation, learning from others and tolerance towards reaching a common goal"),
					app.P().Text("- Enact veto and discuss further if needed"),
					app.P().Text("- Reach a consensus"),
					app.P().Text("- Be a community!"),
				),
			),
		),
		app.Div().ID("main-tutorial-6").Class("pure-u-1 tutorial").Body(
			app.Div().Class("email-content").Body(
				app.Div().Class("email-content-header pure-g").Body(
					app.Div().Class("pure-u-1").Body(
						app.H1().Class("email-content-title").Text("Details"),
					),
				),
				app.Div().Class("email-content-body").Body(
					app.P().Text("- Consensus based, majority democracy leads to a dictatorship of the majority over the minority and doesn't make people cooperate until a compromise is reached"),
					app.P().Text("- No voting, we need to learn how to compromise and reach consensus"),
					app.P().Text("- No moderation and no censorship, we need to be able to overcome hate and spam and still reach a decision"),
				),
			),
		),
	)
}
