package routes

import (
	"net/http"
	"time"
)

func ShowIndexPage(w http.ResponseWriter, r *http.Request) {
	year := time.Now().Year()

	data := map[string]any{
		"Title": "Go & HTMx Demo",
		"Year":  year,
	}

	if err := tmpl.ExecuteTemplate(w, "index", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
