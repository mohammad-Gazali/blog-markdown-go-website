package server

import (
	"net/http"
	"strings"
)

func AddingHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", HomeAndArticleHandler)
}

// handlers
func HomeAndArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		RenderTemplate(w, RenderContext{
			"Title": "Markdown Blog",
			"Articles": GetAllMarkdownFiles(),
		}, "index.html", "partials/card.html")
	} else {
		parts := strings.Split(r.URL.Path, "/")

		slug := parts[len(parts) - 1]

		article := GetMarkdownBySlug(slug)

		if article == nil {
			NotFoundError(w)
			return
		}

		RenderTemplate(w, RenderContext{
			"Title": article.Title,
			"Article": article,
		}, "article.html")
	}
}
