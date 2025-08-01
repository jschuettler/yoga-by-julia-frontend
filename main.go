package main

import (
	"log"

	"github.com/janmarkuslanger/ssgo/builder"
	"github.com/janmarkuslanger/ssgo/page"
	"github.com/janmarkuslanger/ssgo/rendering"
	"github.com/janmarkuslanger/ssgo/writer"
	"github.com/jschuettler/yoga-by-julia-frontend/fetch"
)

func main() {
	renderer := rendering.HTMLRenderer{
		Layout: []string{"templates/layout.html"},
	}

	generator := page.Generator{
		Config: page.Config{
			Pattern:  "/:slug",
			Template: "templates/dynamicpage.html",
			GetPaths: func() []string {
				s, _ := fetch.GetAllPageSlugs()
				return s
			},
			GetData: func(p page.PagePayload) map[string]any {
				return map[string]any{
					"Title":   p.Path,
					"Content": p.Params["slug"],
				}
			},
			Renderer: renderer,
		},
	}

	b := builder.Builder{
		OutputDir: "public",
		Writer:    &writer.FileWriter{},
		Pages: []page.Generator{
			generator,
		},
	}

	if err := b.Build(); err != nil {
		log.Fatal(err)
	}
}
