package pages

import (
	"strings"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// UUIDGen is the main component of the application.
type UUIDGen struct {
	app.Compo
	Version    string
	Quantity   int
	UseHyphens bool
	Uppercase  bool
	Result     string
}

// OnMount initializes the component with default data.
func (p *UUIDGen) OnMount(ctx app.Context) {
	p.Version = "4"
	p.Quantity = 1
	p.UseHyphens = true
	p.Uppercase = false
	p.Generate(ctx)
}

func (p *UUIDGen) OnVersionChange(ctx app.Context, e app.Event) {
	p.Version = ctx.JSSrc().Get("value").String()
}

func (p *UUIDGen) OnQuantityChange(ctx app.Context, e app.Event) {
	val := ctx.JSSrc().Get("value").Int()
	if val < 1 {
		val = 1
	} else if val > 100 {
		val = 100
	}
	p.Quantity = val
}

func (p *UUIDGen) OnHyphensChange(ctx app.Context, e app.Event) {
	p.UseHyphens = ctx.JSSrc().Get("checked").Bool()
}

func (p *UUIDGen) OnUppercaseChange(ctx app.Context, e app.Event) {
	p.Uppercase = ctx.JSSrc().Get("checked").Bool()
}

func (p *UUIDGen) OnGenerate(ctx app.Context, e app.Event) {
	p.Generate(ctx)
}

func (p *UUIDGen) Generate(ctx app.Context) {
	var uuids []string
	for i := 0; i < p.Quantity; i++ {
		var u uuid.UUID
		var err error

		switch p.Version {
		case "1":
			u, err = uuid.NewUUID()
		case "4":
			u, err = uuid.NewRandom()
		case "7":
			u, err = uuid.NewV7()
		default:
			u, err = uuid.NewRandom()
		}

		if err != nil {
			uuids = append(uuids, "Error: "+err.Error())
			continue
		}

		s := u.String()
		if !p.UseHyphens {
			s = strings.ReplaceAll(s, "-", "")
		}
		if p.Uppercase {
			s = strings.ToUpper(s)
		}
		uuids = append(uuids, s)
	}
	p.Result = strings.Join(uuids, "\n")
	ctx.Update()
}

func (p *UUIDGen) OnCopy(ctx app.Context, e app.Event) {
	app.Window().Get("navigator").Get("clipboard").Call("writeText", p.Result)
}

// Render describes the UI.
func (p *UUIDGen) Render() app.UI {
	return app.Div().Class("page-wrapper with-navbar").Body(
		// Navbar
		app.Nav().Class("navbar").Body(
			app.A().Href("#").Class("navbar-brand").Text("Minimalist UUID Generator"),
			app.Ul().Class("navbar-nav d-none d-md-flex").Body(
				app.Li().Class("nav-item").Body(
					app.A().Class("nav-link").Href("https://github.com/google/uuid").Target("_blank").Text("UUID Library"),
				),
			),
		),
		// Content
		app.Div().Class("content-wrapper").Body(
			app.Div().Class("container-fluid").Body(
				app.Div().Class("row").Body(
					app.Div().Class("col-12 col-md-4").Body(
						app.Div().Class("card p-20").Body( // Settings Card
							app.H2().Class("card-title").Text("Settings"),

							// Version Select
							app.Div().Class("form-group").Body(
								app.Label().For("version-select").Text("UUID Version"),
								app.Select().Class("form-control").ID("version-select").OnChange(p.OnVersionChange).Body(
									app.Option().Value("1").Text("Version 1 (Time-based)"),
									app.Option().Value("4").Text("Version 4 (Random)").Selected(true),
									app.Option().Value("7").Text("Version 7 (Time-ordered)"),
								),
							),

							// Quantity Input
							app.Div().Class("form-group").Body(
								app.Label().For("quantity-input").Text("Quantity (1-100)"),
								app.Input().Type("number").Class("form-control").ID("quantity-input").
									Min(1).Max(100).Value(p.Quantity).OnChange(p.OnQuantityChange),
							),

							// Checkboxes
							app.Div().Class("custom-control custom-checkbox mb-10").Body(
								app.Input().Type("checkbox").Class("custom-control-input").ID("hyphens-check").
									Checked(p.UseHyphens).OnChange(p.OnHyphensChange),
								app.Label().Class("custom-control-label").For("hyphens-check").Text("Hyphens"),
							),
							app.Div().Class("custom-control custom-checkbox mb-20").Body(
								app.Input().Type("checkbox").Class("custom-control-input").ID("uppercase-check").
									Checked(p.Uppercase).OnChange(p.OnUppercaseChange),
								app.Label().Class("custom-control-label").For("uppercase-check").Text("Uppercase"),
							),

							// Generate Button
							app.Button().Class("btn btn-primary btn-block").OnClick(p.OnGenerate).Text("Generate UUIDs"),
						),
					),
					app.Div().Class("col-12 col-md-8").Body(
						app.Div().Class("card p-20").Body( // Result Card
							app.Div().Class("d-flex justify-content-between align-items-center mb-10").Body(
								app.H2().Class("card-title m-0").Text("Result"),
								app.Button().Class("btn").OnClick(p.OnCopy).Text("Copy to Clipboard"),
							),
							app.Textarea().Class("form-control").Style("height", "400px").Style("font-family", "monospace").
								ReadOnly(true).Body(app.Text(p.Result)),
						),
					),
				),
			),
		),
	)
}
