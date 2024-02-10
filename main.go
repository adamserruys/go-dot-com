package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type handlerFunc func(c *Context) error

type Context struct {
	w   http.ResponseWriter
	req *http.Request
}

func (c *Context) render(pageTmpl string, data any) error {
	baseTmpl := "./src/templates/base.html"

	tmpls := []string{baseTmpl, pageTmpl}

	ts, err := template.ParseFiles(tmpls...)
	if err != nil {
		return err
	}

	return ts.ExecuteTemplate(c.w, "base", data)
}

func (c *Context) writeJSON(code int, v any) error {
	c.w.Header().Set("Content-Type", "application/json")
	c.w.WriteHeader(code)
	return json.NewEncoder(c.w).Encode(v)
}

func main() {
	port := ":42069"
	styles := http.FileServer(http.Dir("./src/assets/styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", styles))

	http.HandleFunc("/", makeHandlerFunc(handleHome))
	http.HandleFunc("/about", makeHandlerFunc(handleAbout))
	http.HandleFunc("/music", makeHandlerFunc(handleMusic))

	fmt.Println("Listening on port: ", port)
	http.ListenAndServe(port, nil)
}

func makeHandlerFunc(fn handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{
			w:   w,
			req: r,
		}
		if err := fn(ctx); err != nil {
			ctx.writeJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}

func handleHome(c *Context) error {
	return c.render("./src/templates/index.html", nil)
}

func handleAbout(c *Context) error {
	return c.render("./src/templates/about.html", nil)
}

func handleMusic(c *Context) error {
	return c.render("./src/templates/music.html", nil)
}
