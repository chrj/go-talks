// Package talks renders the present files of this repository as static pages.
//
// The parser and the element templates come from golang.org/x/tools/present.
// The page templates are in this package, because the upstream templates
// belong to the present command, which serves the slides from a web server.
package talks

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"

	"golang.org/x/tools/present"
)

//go:embed templates/*.tmpl
var templates embed.FS

// Talk describes one presentation.
type Talk struct {
	// Source is the path of the .slide file, relative to the source
	// directory, with forward slashes.
	Source string
	// Page is the path of the rendered page, relative to the output
	// directory, with forward slashes.
	Page string
	// Group is the first element of Source, which is the year.
	Group    string
	Title    string
	Subtitle string
}

// Group holds the talks that share the first element of their path.
type Group struct {
	Name  string
	Talks []Talk
}

// IndexPage holds the text of the landing page that does not come from the
// talks.
type IndexPage struct {
	Title string
	// Website is the address of the main site, which the landing page links
	// back to. WebsiteLabel is the text of that link.
	Website      string
	WebsiteLabel string
	Repository   string
}

// Render reads the present file at fsPath and writes its page to w. The
// present parser reads the .code and .play files against the directory of
// fsPath. relPath is the same file as a path relative to the source directory,
// and it gives the page its address on the site.
//
// Render is not safe for concurrent use. The present package keeps the
// playground switch in a package level variable, which Render sets.
func Render(w io.Writer, fsPath, relPath string) (Talk, error) {
	// The present parser reads this variable when it parses a .play
	// directive. Without it, the snippets get no Run button.
	present.PlayEnabled = true

	b, err := os.ReadFile(fsPath)
	if err != nil {
		return Talk{}, fmt.Errorf("read present file: %w", err)
	}

	doc, err := present.Parse(bytes.NewReader(b), fsPath, 0)
	if err != nil {
		return Talk{}, fmt.Errorf("parse present file %s: %w", fsPath, err)
	}

	tmpl, err := slideTemplate()
	if err != nil {
		return Talk{}, err
	}

	if err := doc.Render(w, tmpl); err != nil {
		return Talk{}, fmt.Errorf("render present file %s: %w", fsPath, err)
	}

	return Talk{
		Source:   relPath,
		Page:     pagePath(relPath),
		Group:    groupName(relPath),
		Title:    doc.Title,
		Subtitle: doc.Subtitle,
	}, nil
}

// Index writes the landing page, which links to each talk.
func Index(w io.Writer, page IndexPage, talks []Talk) error {
	tmpl, err := template.ParseFS(templates, "templates/index.tmpl")
	if err != nil {
		return fmt.Errorf("parse the index template: %w", err)
	}

	data := struct {
		IndexPage
		Groups []Group
	}{page, Groups(talks)}

	if err := tmpl.ExecuteTemplate(w, "index", data); err != nil {
		return fmt.Errorf("render the index page: %w", err)
	}
	return nil
}

// Groups sorts the talks into one group for each year. The newest group comes
// first, and the talks inside a group are in the order of their path.
func Groups(talks []Talk) []Group {
	byName := make(map[string][]Talk)
	for _, t := range talks {
		byName[t.Group] = append(byName[t.Group], t)
	}

	groups := make([]Group, 0, len(byName))
	for name, ts := range byName {
		slices.SortFunc(ts, func(a, b Talk) int { return strings.Compare(a.Page, b.Page) })
		groups = append(groups, Group{Name: name, Talks: ts})
	}
	slices.SortFunc(groups, func(a, b Group) int { return strings.Compare(b.Name, a.Name) })
	return groups
}

// StaticFiles gives the files that the render command copies from the static
// directory of the present command. The pages load styles.css through
// slides.js, which builds the address from the /static/ prefix.
func StaticFiles() []string {
	return []string{"favicon.ico", "slides.js", "styles.css"}
}

// playScripts are the playground scripts of the present command, in the order
// that the pages load them.
var playScripts = []string{"jquery.js", "jquery-ui.js", "playground.js", "play.js"}

// PlayScript joins the playground scripts from the static directory of the
// present command into one file, and appends the call that starts the
// playground. The HTTP transport posts each snippet to /compile on the same
// origin, which the Worker proxies to the Go playground.
func PlayScript(static fs.FS) ([]byte, error) {
	var buf bytes.Buffer
	for _, name := range playScripts {
		b, err := fs.ReadFile(static, name)
		if err != nil {
			return nil, fmt.Errorf("read the playground script %s: %w", name, err)
		}
		buf.Write(b)
		buf.WriteString("\n")
	}
	buf.WriteString("initPlayground(new HTTPTransport());\n")
	return buf.Bytes(), nil
}

// slideTemplate parses the templates that render one .slide file. The element
// templates of the present package need the playable function.
func slideTemplate() (*template.Template, error) {
	tmpl := present.Template().Funcs(template.FuncMap{"playable": playable})
	tmpl, err := tmpl.ParseFS(templates, "templates/action.tmpl", "templates/slides.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parse the slide templates: %w", err)
	}
	return tmpl, nil
}

// playable reports whether the page gives a snippet a Run button. The
// playground at play.golang.org runs Go files only.
func playable(c present.Code) bool {
	return c.Play && c.Ext == ".go"
}

// pagePath gives the address of the rendered page for a source path.
func pagePath(relPath string) string {
	return strings.TrimSuffix(relPath, path.Ext(relPath)) + ".html"
}

// groupName gives the first element of a path, or an empty string when the
// path has one element only.
func groupName(relPath string) string {
	dir, _ := path.Split(relPath)
	if dir == "" {
		return ""
	}
	return strings.SplitN(strings.TrimSuffix(dir, "/"), "/", 2)[0]
}
