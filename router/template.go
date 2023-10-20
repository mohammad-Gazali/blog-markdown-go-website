package router

import (
	"html/template"
	"net/http"
)

// this function return the html template file as a byte array
func MustGetTemplate(w http.ResponseWriter, fileName string, context any) {
	tmpl, err := template.ParseFiles("./views/" + fileName)

	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-Type", "application/html")

	err = tmpl.Execute(w, context)

	if err != nil {
		panic(err)
	}
}