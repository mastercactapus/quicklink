.PHONY: clean start

build: build/bin/quicklink
clean:
	rm -rf build

start: build/bin/quicklink
	@if [ ! -f build/dev.txt ]; then echo "example -> http://example.com" > build/dev.txt; fi
	go run . -addr localhost:8080 -txt build/dev.txt

build/tools/sqlc: go.mod go.sum
	@echo "  >  Building sqlc"
	@go build -o $@ github.com/kyleconroy/sqlc/cmd/sqlc

pkg/store/pgstore/pg-query.sql.go: pkg/store/pg-query.sql pkg/store/pg-schema.sql sqlc.yaml build/tools/sqlc
	@echo "  >  Generating pg-query.sql.go"
	@./build/tools/sqlc generate

build/bin/quicklink: Makefile pkg/store/pgstore/pg-query.sql.go $(find . -name '*.go') internal/web/index.html go.mod go.sum
	@echo "  >  Building quicklink into $@"
	@mkdir -p build/bin
	@CGO_ENABLED=0 go build -trimpath -o $@
