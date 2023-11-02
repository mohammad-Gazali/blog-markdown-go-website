package server

import (
	"html/template"
	"net/http"
)

// this function return the html template file as a byte array
func RenderTemplate(w http.ResponseWriter, context any, files ...string) {

	for i, f := range files {
		files[i] = "./views/" + f
	}

	files = append(files,
		"./views/base.html",
		"./views/partials/navbar.html",
		"./views/partials/footer.html",
	)

	tmpl, err := template.ParseFiles(files...)

	if err != nil {
		InternalServerError(w)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base.html", context)

	if err != nil {
		InternalServerError(w)
		return
	}
}
