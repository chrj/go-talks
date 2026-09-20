// The site is static. The one exception is the Run button on the slides: the
// playground scripts post the snippet to /compile on this origin, because the
// Go playground sends no CORS headers. This Worker sends those requests on.
//
// Cloudflare serves the files in public/ before it calls this Worker, so only
// the paths below and the paths that do not exist arrive here.

const PLAYGROUND = "https://play.golang.org";

// The paths that the playground scripts call. /compile runs a snippet. /vet is
// the fallback for an older playground that does not vet in the same request.
const PROXIED = new Set(["/compile", "/vet"]);

export default {
  async fetch(request, env) {
    const { pathname } = new URL(request.url);

    if (!PROXIED.has(pathname)) {
      return env.ASSETS.fetch(request);
    }

    if (request.method !== "POST") {
      return new Response(`The path ${pathname} takes POST requests only.\n`, {
        status: 405,
        headers: { "Allow": "POST", "Content-Type": "text/plain; charset=utf-8" },
      });
    }

    try {
      return await fetch(PLAYGROUND + pathname, {
        method: "POST",
        headers: {
          "Content-Type":
            request.headers.get("Content-Type") ?? "application/x-www-form-urlencoded",
        },
        body: await request.text(),
      });
    } catch (err) {
      console.error(`POST ${PLAYGROUND}${pathname} failed: ${err.stack ?? err}`);
      return new Response(
        `The Go playground at ${PLAYGROUND} did not answer: ${err.message}\n` +
          `Run the snippet on https://go.dev/play instead.\n`,
        { status: 502, headers: { "Content-Type": "text/plain; charset=utf-8" } },
      );
    }
  },
};
