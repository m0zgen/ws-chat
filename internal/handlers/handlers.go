package handlers

import (
	"github.com/CloudyKit/jet/v6"
	"net/http"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode())

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return err
	}

	return nil

}
