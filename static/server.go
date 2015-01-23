package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

const LOCAL_URL string = "/static/"
const LOCAL_ROOT string = "assets/"

type Context struct {
	Title  string
	Static string
}

func Home(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Welcome!"}
	render(w, "index", context)
}

func About(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "About"}
	render(w, "about", context)
}

func render(w http.ResponseWriter, tmpl string, context Context) {
	context.Static = LOCAL_URL
	tmpl_list := []string{"templates/base.html",
		fmt.Sprintf("templates/%s.html", tmpl)}
	t, err := template.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, context)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	static_file := req.URL.Path[len(LOCAL_URL):]
	if len(static_file) != 0 {
		f, err := http.Dir(LOCAL_ROOT).Open(static_file)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		}
	}
	http.NotFound(w, req)
}

func main() {
	log.Println("starting server on http://localhost:8000/")
	http.HandleFunc("/", Home)
	http.HandleFunc("/about/", About)
	http.HandleFunc(LOCAL_URL, StaticHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
