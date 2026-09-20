# go-talks

The talks I have given on Go, at [talks.technobabble.dk](https://talks.technobabble.dk).

## 2017

- [Introduction to Go](https://talks.technobabble.dk/2017/intro/presentation)
- [The Go X repository](https://talks.technobabble.dk/2017/x-packages/presentation)

## How the site is built

Each talk is a `.slide` file in the present format. The `render` command reads
these files with `golang.org/x/tools/present` and writes one HTML page for each
of them into `public/`. It also writes the landing page and copies the
stylesheet and the scripts that the pages need.

The static files belong to the present command, which is part of
`golang.org/x/tools`. The render command finds them with `go list`, so that
they always match the module version in `go.mod`. Nothing from the present
command is a copy in this repository, except the two templates in
`internal/talks/templates/`.

`public/` is generated, so Git ignores it. GitHub Actions builds the site again
on each push to `master`.

## The playground

The slides of the introduction talk hold 10 runnable snippets. The scripts of
the playground post each snippet to `/compile` on the same origin, because the
Go playground sends no CORS headers. The Worker script in `src/index.js` sends
those requests to `https://play.golang.org`. Every other request gets a static
file.

Only Go snippets get a Run button. The playground cannot run a shell script.

## Deploy the site

The site is a Cloudflare Worker. `wrangler.jsonc` holds the configuration of
the Worker and its custom domain. The `deploy.yml` workflow renders the talks,
copies `_headers` into `public/`, and then runs `wrangler deploy`.

`_headers` gives the HTTP headers that Cloudflare adds to the responses, for
example HSTS.

For each pull request, the `build.yml` workflow uploads a preview version of
the Worker. The preview version does not change the deployed site. It is on
`https://pr-<number>-go-talks.<account>.workers.dev`, and the pull request
links to it with "View deployment". Pull requests from forks and from
Dependabot get no preview, because they get no repository secrets.

A version needs a Worker, and the first deploy on `master` makes it. Until that
deploy, the preview step writes a notice and stops.

The workflows need two repository secrets:

- `CLOUDFLARE_API_TOKEN` — an API token with the "Edit Cloudflare Workers"
  template
- `CLOUDFLARE_ACCOUNT_ID` — the ID of the Cloudflare account

## Build the site on your machine

    go run ./cmd/render

The command writes `public/`. It replaces the files that are already
there, but it removes nothing. If you rename or delete a talk, delete `public/` first.

To see the site, run the Worker with the local Cloudflare runtime:

    cp _headers public/
    npx wrangler@4 dev

## Run the tests

    go test -race ./...

    golangci-lint run ./...

## Security scanning

Dependabot opens a pull request each week for the GitHub Actions and the Go
modules. The `security-daily.yml` workflow runs govulncheck, Semgrep, and gosec
each night, and it also runs on demand.

Both Semgrep and gosec leave out `2017/`. The files there are the code of the
talks. They teach a point, and some of them show a password or a server without
a timeout. A scan of them reports the talk, not the site. The `ignore`
directive in `go.mod` is not enough, because these two tools read the files
themselves and not the build list.

To run the same scans on your machine:

    go install golang.org/x/vuln/cmd/govulncheck@v1.7.0
    govulncheck ./...

    go install github.com/securego/gosec/v2/cmd/gosec@v2.28.0
    gosec -exclude-generated -exclude-dir=2017 ./...

## Add a talk

1. Make a directory for the talk, under the year: `2026/my-talk/`.
2. Write the presentation in `presentation.slide`, in the present format.
3. Put the code of each snippet in its own file, next to the slide file.
4. Open a pull request. The preview version shows the talk before the deploy.

The landing page groups the talks by the first element of their path, and it
puts the newest year first.

Note: the `ignore` directive in `go.mod` makes the go command ignore the
snippets. They are examples, not packages, and several hold their own
`func main` in one directory.
