# flies

A simple debugging webserver which produces detailed logs of each request it receives.

## Usage

Compile, run, and follow logs:

```bash
go run .
```

Optionally, set the `FLIES_FORMAT` environment variable to one of the following
formats:

- `pretty` (default): Pretty print the details of the request.
- `json`: Print one line of JSON per request (useful when piped to `jq`).
- `wire`: Print request in wire format.
- `template`: Print request according to a go html/template.

```bash
FLIES_FORMAT=json go run .
```

See other environment variables in [the CLI docs](cli/doc.go).
