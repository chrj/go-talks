module github.com/chrj/go-talks

go 1.27.1

// The files of the talks are snippets, not packages. Several hold their own
// func main in one directory, and some import modules that no longer exist.
// The go command must not try to build them.
ignore ./2017

require golang.org/x/tools v0.50.0

require github.com/yuin/goldmark v1.7.17 // indirect
