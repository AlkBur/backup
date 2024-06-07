package routes

import (
	"html/template"
)

var tmpl *template.Template

func init() {
	if tmpl == nil {
		if tmpl == nil {
			tmpl = template.Must(tmpl.ParseGlob("views/layouts/*.html"))
			template.Must(tmpl.ParseGlob("views/*.html"))
		}
	}
}
