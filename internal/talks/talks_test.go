package talks_test

import (
	"bytes"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/chrj/go-talks/internal/talks"
)

func TestRender(t *testing.T) {
	var buf bytes.Buffer

	talk, err := talks.Render(&buf, "testdata/example/talk.slide", "2017/example/talk.slide")
	if err != nil {
		t.Fatalf("Render() error = %v, want no error", err)
	}

	want := talks.Talk{
		Source:   "2017/example/talk.slide",
		Page:     "2017/example/talk.html",
		Group:    "2017",
		Title:    "Example talk",
		Subtitle: "A subtitle for the test",
	}
	if talk != want {
		t.Errorf("Render() talk = %+v, want %+v", talk, want)
	}

	page := buf.String()
	for _, want := range []string{
		"<title>Example talk</title>",
		"<h1>Example talk</h1>",
		"<h3>A subtitle for the test</h3>",
		"<h3>A section</h3>",
		"first bullet",
		`fmt.Println(&#34;Hello, talk&#34;)`,
		`<script src="/static/play.js"></script>`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("Render() page does not hold %q", want)
		}
	}
}

func TestRenderGivesRunButtonToGoSnippetsOnly(t *testing.T) {
	var buf bytes.Buffer

	if _, err := talks.Render(&buf, "testdata/example/talk.slide", "2017/example/talk.slide"); err != nil {
		t.Fatalf("Render() error = %v, want no error", err)
	}

	page := buf.String()
	if got := strings.Count(page, `class="code playground"`); got != 1 {
		t.Errorf("playground snippets = %d, want 1 (the Go file only)", got)
	}
	if !strings.Contains(page, "go build ./...") {
		t.Error("Render() page does not hold the shell snippet")
	}
}

func TestRenderErrors(t *testing.T) {
	tests := []struct {
		name    string
		fsPath  string
		wantMsg string
	}{
		{
			name:    "file does not exist",
			fsPath:  "testdata/no-such-talk.slide",
			wantMsg: "read present file",
		},
		{
			name:    "code directive names a file that does not exist",
			fsPath:  "testdata/broken.slide",
			wantMsg: "parse present file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			_, err := talks.Render(&buf, tt.fsPath, "talk.slide")
			if err == nil {
				t.Fatalf("Render() error = nil, want an error that holds %q", tt.wantMsg)
			}
			if !strings.Contains(err.Error(), tt.wantMsg) {
				t.Errorf("Render() error = %q, want it to hold %q", err, tt.wantMsg)
			}
		})
	}
}

func TestGroups(t *testing.T) {
	tests := []struct {
		name  string
		talks []talks.Talk
		want  []talks.Group
	}{
		{
			name:  "no talks",
			talks: nil,
			want:  []talks.Group{},
		},
		{
			name: "newest group first, talks in path order",
			talks: []talks.Talk{
				{Page: "2017/intro.html", Group: "2017", Title: "Intro"},
				{Page: "2019/b.html", Group: "2019", Title: "B"},
				{Page: "2017/a.html", Group: "2017", Title: "A"},
				{Page: "2019/a.html", Group: "2019", Title: "A"},
			},
			want: []talks.Group{
				{Name: "2019", Talks: []talks.Talk{
					{Page: "2019/a.html", Group: "2019", Title: "A"},
					{Page: "2019/b.html", Group: "2019", Title: "B"},
				}},
				{Name: "2017", Talks: []talks.Talk{
					{Page: "2017/a.html", Group: "2017", Title: "A"},
					{Page: "2017/intro.html", Group: "2017", Title: "Intro"},
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := talks.Groups(tt.talks)

			if len(got) != len(tt.want) {
				t.Fatalf("Groups() = %+v, want %+v", got, tt.want)
			}
			for i := range got {
				if got[i].Name != tt.want[i].Name {
					t.Errorf("Groups()[%d].Name = %q, want %q", i, got[i].Name, tt.want[i].Name)
				}
				if len(got[i].Talks) != len(tt.want[i].Talks) {
					t.Fatalf("Groups()[%d].Talks = %+v, want %+v", i, got[i].Talks, tt.want[i].Talks)
				}
				for j := range got[i].Talks {
					if got[i].Talks[j] != tt.want[i].Talks[j] {
						t.Errorf("Groups()[%d].Talks[%d] = %+v, want %+v", i, j, got[i].Talks[j], tt.want[i].Talks[j])
					}
				}
			}
		})
	}
}

func TestIndex(t *testing.T) {
	var buf bytes.Buffer

	page := talks.IndexPage{Title: "Go Talks", Repository: "https://github.com/chrj/go-talks"}
	list := []talks.Talk{
		{Page: "2017/intro/presentation.html", Group: "2017", Title: "Introduction to Go", Subtitle: "The not so short version"},
		{Page: "2019/x/presentation.html", Group: "2019", Title: "The Go X repository"},
	}

	if err := talks.Index(&buf, page, list); err != nil {
		t.Fatalf("Index() error = %v, want no error", err)
	}

	out := buf.String()
	for _, want := range []string{
		"<title>Go Talks</title>",
		`<a href="/2017/intro/presentation.html">Introduction to Go</a>`,
		`<span class="subtitle">The not so short version</span>`,
		`<a href="/2019/x/presentation.html">The Go X repository</a>`,
		`<a href="https://github.com/chrj/go-talks">Source on GitHub</a>`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Index() page does not hold %q", want)
		}
	}

	if newest, oldest := strings.Index(out, ">2019<"), strings.Index(out, ">2017<"); newest > oldest {
		t.Errorf("Index() puts 2017 before 2019, want the newest year first")
	}
}

func TestPlayScript(t *testing.T) {
	static := fstest.MapFS{
		"jquery.js":     {Data: []byte("jquery")},
		"jquery-ui.js":  {Data: []byte("jquery-ui")},
		"playground.js": {Data: []byte("playground")},
		"play.js":       {Data: []byte("play")},
	}

	got, err := talks.PlayScript(static)
	if err != nil {
		t.Fatalf("PlayScript() error = %v, want no error", err)
	}

	want := "jquery\njquery-ui\nplayground\nplay\ninitPlayground(new HTTPTransport());\n"
	if string(got) != want {
		t.Errorf("PlayScript() = %q, want %q", got, want)
	}
}

func TestPlayScriptMissingScript(t *testing.T) {
	static := fstest.MapFS{"jquery.js": {Data: []byte("jquery")}}

	_, err := talks.PlayScript(static)
	if err == nil {
		t.Fatal("PlayScript() error = nil, want an error that names the missing script")
	}
	if !strings.Contains(err.Error(), "jquery-ui.js") {
		t.Errorf("PlayScript() error = %q, want it to name jquery-ui.js", err)
	}
}
