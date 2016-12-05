// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command blog is a web server for the Go blog that can run on App Engine or
// as a stand-alone HTTP server.
package main

import (
	"errors"
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"golang.org/x/tools/blog"
	"golang.org/x/tools/present"
	"golang.org/x/tools/godoc/static"

	_ "golang.org/x/tools/playground"
)

const hostname = "blog.golang.org" // default hostname for blog server

var (
	httpAddr = flag.String("http", "0.0.0.0:8080", "HTTP listen address")
	barContentPath = flag.String("bar", "bar/", "path to bar content files")
	contentPath = flag.String("content", "event/", "path to content files")
	templatePath = flag.String("template", "template/", "path to template files")
	staticPath = flag.String("static", "static/", "path to static files")

	config = blog.Config{
		Hostname:     hostname,
		BaseURL:      "//" + hostname,
		GodocURL:     "//golang.org",
		HomeArticles: 5, // articles to display on the home page
		FeedArticles: 10, // articles to include in Atom and JSON feeds
		PlayEnabled:  true,
		FeedTitle:    "Weekend Chinese Tech Meetup",
	}
)

func init() {
	present.Register("md", parseMarkdown)
}

type Markdown struct {
	template.HTML
}

func (i Markdown) TemplateName() string { return "md" }

func parseMarkdown(ctx *present.Context, fileName string, lineno int, text string) (present.Elem, error) {
	p := strings.Fields(text)
	if len(p) != 2 {
		return nil, errors.New("invalid .md args")
	}
	name := filepath.Join(filepath.Dir(fileName), p[1])
	b, err := ctx.ReadFile(name)
	if err != nil {
		return nil, err
	}

	unsafe := blackfriday.MarkdownCommon(b)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	// TODO(lamuguo): don't use "present.HTML" here.
	return Markdown{template.HTML(html)}, nil
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	b, ok := static.Files[name]
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, name, time.Time{}, strings.NewReader(b))
}

func serveBlog(prefix string, contentPath string) {
	blogConfig := blog.Config{
		Hostname: hostname,
		BaseURL: "//" + hostname + prefix,
		HomeArticles: 5,
		FeedArticles: 10,
		PlayEnabled: true,
		FeedTitle: "TechM",
		ContentPath: contentPath,
		TemplatePath: *templatePath,
	}

	s2, err := blog.NewServer(blogConfig)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle(prefix + "/", http.StripPrefix(prefix, s2))
}

func redirectToBarPath(path string) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/bar" + path, http.StatusFound)
	})
}

func main() {
	flag.Parse()

	http.Handle("/lib/godoc/", http.StripPrefix("/lib/godoc/", http.HandlerFunc(staticHandler)))

	serveBlog("/bar", *barContentPath)
	serveBlog("", *contentPath)

	redirectToBarPath("/about")
	redirectToBarPath("/groups")
	redirectToBarPath("/members")
	redirectToBarPath("/wechat")


	log.Printf("xfguo: hello world\n")
	//
	//redirect := func(w http.ResponseWriter, r *http.Request) {
	//	http.Redirect(w, r, "/event/", http.StatusFound)
	//}
	//http.HandleFunc("/", redirect)
	//http.HandleFunc("/blog/", redirect)

	fs := http.FileServer(http.Dir(*staticPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
