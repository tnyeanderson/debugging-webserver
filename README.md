# flies

A simple debugging webserver which produces detailed logs of each request it receives.

## Usage

Compile, run, and follow logs:

```bash
go run .
```

Optionally, set the `FLIES_LOG_FORMAT` environment variable to one of the
following formats:

- `text` (default): Pretty print the details of the request
- `json`: Print one line of JSON per request (useful when piped to `jq`)

```bash
FLIES_LOG_FORMAT=json go run .
```

