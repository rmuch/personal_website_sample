package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Page view model.
type Page struct {
	Title      string // Page title.
	Body       string // Page "body" - unused except for error pages.
	CurrentURL string // Current request URL
}

var funcMap = template.FuncMap{
	"menuLink": func(target string, label string, current string) template.HTML {
		// Quick workaround - request to root will appear as an empty string.
		if current == "" {
			current = "/"
		}

		// target, link, current
		if current == target {
			return template.HTML("<a class=\"active\" href=\"" + target + "\">" + label + "</a>")
		} else {
			return template.HTML("<a href=\"" + target + "\">" + label + "</a>")
		}
	},
	"getYear": func() string {
		return fmt.Sprintf("%d", time.Now().Year())
	},
}

// Projects page template.
var projectsTmpl = template.Must(template.New("projects").Funcs(funcMap).ParseFiles(
	tmplRoot+"_base.html",
	tmplRoot+"_menu.html",
	tmplRoot+"projects.html"))

// GET /projects
func projectsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		errNotFound(w, r)
		return
	}

	pageModel := &Page{
		Title:      "Projects",
		CurrentURL: r.URL.Path,
	}

	err := projectsTmpl.ExecuteTemplate(w, "base", pageModel)
	if err != nil {
		log.Print(err)

		errInternal(w, r)
	}
}

// About page template.
var aboutTmpl = template.Must(template.New("about").Funcs(funcMap).ParseFiles(
	tmplRoot+"_base.html",
	tmplRoot+"_menu.html",
	tmplRoot+"about.html"))

// GET /about
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		errNotFound(w, r)
		return
	}

	pageModel := &Page{
		Title:      "About",
		CurrentURL: r.URL.Path,
	}

	err := aboutTmpl.ExecuteTemplate(w, "base", pageModel)
	if err != nil {
		log.Print(err)

		errInternal(w, r)
	}
}

// Index page template.
var indexTmpl = template.Must(template.New("index").Funcs(funcMap).ParseFiles(
	tmplRoot+"_base.html",
	tmplRoot+"_menu.html",
	tmplRoot+"index.html"))

// GET /
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Make sure this is a GET.
	if r.Method != "GET" {
		errNotFound(w, r)
		return
	}

	pageModel := &Page{
		Title: "Home",
	}

	err := indexTmpl.ExecuteTemplate(w, "base", pageModel)
	if err != nil {
		log.Print(err)

		errInternal(w, r)
	}
}

// * /
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		homeHandler(w, r)
	} else {
		log.Printf("rootHandler: Could not forward request for %s any further.", r.RequestURI)

		errNotFound(w, r)
	}
}

// Error page template.
var errorTmpl = template.Must(template.New("error").Funcs(funcMap).ParseFiles(
	tmplRoot+"_base.html",
	tmplRoot+"_menu.html",
	tmplRoot+"error.html"))

// Error 404
func errNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)

	pageModel := &Page{
		Title: "404 Not Found",
		Body:  "The page you were looking for was not found.",
	}

	err := errorTmpl.ExecuteTemplate(w, "base", pageModel)
	if err != nil {
		log.Print(err)

		// TODO: Return internal server error.
		fmt.Fprint(w, "<h1>404 Not Found</h1>")
	}
}

// Error 500
func errInternal(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)

	pageModel := &Page{
		Title: "500 Internal Server Error",
		Body:  "An error occurred while attempting to load the page.",
	}

	err := errorTmpl.ExecuteTemplate(w, "base", pageModel)
	if err != nil {
		log.Print(err)

		// TODO: Return internal server error.
		fmt.Fprint(w, "<h1>500 Internal Server Error</h1>")
	}
}
