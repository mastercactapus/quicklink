# QuickLink

An ultra-simple redirect service.

## Usage

The root URL (default is http://localhost:8080) has a form that lets you add a link and a destination.

For example, if you add a link to `google` and a destination of `https://google.com`, then you can access `http://localhost:8080/google` and it will redirect you to `https://google.com`.

### Local Service

Ideally, it should run on port 80 on a short hostname (e.g., `ql` as in `ql/mylink`), but you can run it on any port and any hostname.

## Installation

Using `go install`:

```bash
go install github.com/mastercactapus/quicklink@latest
```

By default, it will use an in-memory data store. You can also use a `.txt` file, or a Postgres DB (it will create its table the first time it connects).

```bash
$ > quicklink -h
Usage of quicklink:
  -addr string
    	http service address (default ":8080")
  -pg string
    	postgres connection string
  -txt string
    	text file to use for persistence
```

## Development

Requires Go >= 1.19.

To start the server for development run `make start`
To build run `make`, the binary will be in `./build/bin/quicklink`
