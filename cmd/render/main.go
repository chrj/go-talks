// Command render builds the static site of this repository.
//
// It reads each .slide file under the source directory and writes one HTML
// page for it in the output directory. Then it writes the landing page and
// copies the static files that the pages need from the present command.
//
// The static files belong to golang.org/x/tools, which go.mod pins. The
// command finds them with "go list", so that they always match the module in
// the build. Give the -present flag to use a different copy.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chrj/go-talks/internal/talks"
)

const (
	siteTitle    = "Go Talks"
	website      = "https://technobabble.dk"
	websiteLabel = "technobabble.dk"
	repository   = "https://github.com/chrj/go-talks"
	toolsModule  = "golang.org/x/tools"
)

func main() {
	log.SetFlags(0)

	src := flag.String("src", ".", "directory that holds the .slide files")
	out := flag.String("out", "public", "directory to write the site to")
	present := flag.String("present", "", "directory of the present command, which holds static/ (default: from "+toolsModule+" in go.mod)")
	flag.Parse()

	if err := run(*src, *out, *present); err != nil {
		log.Fatalf("render: %v", err)
	}
}

func run(src, out, present string) error {
	if present == "" {
		dir, err := findPresentDir()
		if err != nil {
			return err
		}
		present = dir
	}

	sources, err := findSlides(src, out)
	if err != nil {
		return err
	}
	if len(sources) == 0 {
		return fmt.Errorf("found no .slide file under %s: give the -src flag the directory of the talks", src)
	}

	rendered := make([]talks.Talk, 0, len(sources))
	for _, relPath := range sources {
		talk, err := renderTalk(filepath.Join(src, filepath.FromSlash(relPath)), relPath, out)
		if err != nil {
			return err
		}
		rendered = append(rendered, talk)
		log.Printf("rendered %s as %s: %q", relPath, talk.Page, talk.Title)
	}

	if err := writeIndex(out, rendered); err != nil {
		return err
	}
	if err := copyStatic(present, out); err != nil {
		return err
	}

	log.Printf("wrote %d talks to %s, with the static files from %s", len(rendered), out, present)
	return nil
}

// renderTalk renders one present file into memory and then writes the page. A
// parse error leaves no half written page in the output directory.
func renderTalk(fsPath, relPath, out string) (talks.Talk, error) {
	var buf bytes.Buffer
	talk, err := talks.Render(&buf, fsPath, relPath)
	if err != nil {
		return talks.Talk{}, err
	}
	if err := writeFile(filepath.Join(out, filepath.FromSlash(talk.Page)), buf.Bytes()); err != nil {
		return talks.Talk{}, err
	}
	return talk, nil
}

func writeIndex(out string, rendered []talks.Talk) error {
	var buf bytes.Buffer
	page := talks.IndexPage{
		Title:        siteTitle,
		Website:      website,
		WebsiteLabel: websiteLabel,
		Repository:   repository,
	}
	if err := talks.Index(&buf, page, rendered); err != nil {
		return err
	}
	return writeFile(filepath.Join(out, "index.html"), buf.Bytes())
}

// copyStatic writes the files that the pages load from /static/: the icon and
// the slide script as they are, and the stylesheet and the playground script
// with the parts that this repository adds.
func copyStatic(present, out string) error {
	dir := filepath.Join(present, "static")
	static := os.DirFS(dir)

	for _, name := range talks.StaticFiles() {
		b, err := fs.ReadFile(static, name)
		if err != nil {
			return fmt.Errorf("read the static file %s from %s: %w", name, dir, err)
		}
		if err := writeFile(filepath.Join(out, "static", name), b); err != nil {
			return err
		}
	}

	style, err := talks.Stylesheet(static)
	if err != nil {
		return fmt.Errorf("build the stylesheet from %s: %w", dir, err)
	}
	if err := writeFile(filepath.Join(out, "static", talks.StylesheetName), style); err != nil {
		return err
	}

	play, err := talks.PlayScript(static)
	if err != nil {
		return fmt.Errorf("build the playground script from %s: %w", dir, err)
	}
	return writeFile(filepath.Join(out, "static", "play.js"), play)
}

// findSlides gives the .slide files under src as paths with forward slashes,
// relative to src and in lexical order. It skips the output directory, the
// directories that start with a dot, and the testdata directories, which hold
// the fixtures of the tests.
func findSlides(src, out string) ([]string, error) {
	skip, err := filepath.Abs(out)
	if err != nil {
		return nil, fmt.Errorf("resolve the output directory %s: %w", out, err)
	}

	var found []string
	err = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk %s: %w", path, err)
		}
		if d.IsDir() {
			abs, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("resolve %s: %w", path, err)
			}
			if abs == skip || (skipDir(d.Name()) && path != src) {
				return fs.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".slide" {
			return nil
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("resolve %s against %s: %w", path, src, err)
		}
		found = append(found, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return found, nil
}

// skipDir reports whether the walk leaves out a directory and its children.
func skipDir(name string) bool {
	return strings.HasPrefix(name, ".") || name == "testdata"
}

// findPresentDir asks the go command where the module of the present command
// is. The module must be in the build list of go.mod and in the module cache.
func findPresentDir() (string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", toolsModule)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ask the go command for the directory of %s: %w: %s", toolsModule, err, strings.TrimSpace(stderr.String()))
	}

	dir := strings.TrimSpace(string(out))
	if dir == "" {
		return "", fmt.Errorf("the go command gave no directory for %s: run \"go mod download %s\" first", toolsModule, toolsModule)
	}
	return filepath.Join(dir, "cmd", "present"), nil
}

func writeFile(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("make the directory for %s: %w", path, err)
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
