REFLEX_CMD := reflex -d none -s -r "(\.go$$|^Makefile$$|^go.mod$$|^go.sum$$)" --

# Don't treat those as files
.PHONY: upgrade clean distclean build test test-watch prepare

# Disable implicit rules
.SUFFIXES:

all: build test

build: temp/setup-dev
	@go build ./...

test: temp/setup-dev
	@go test -timeout 5s -shuffle on -coverprofile=temp/coverage.out \
			-ldflags="-X \"github.com/akshaal/akgoli/appinfo.version=3.2.1\" -X \"github.com/akshaal/akgoli/appinfo.idName=test\"" \
			-cover ./... \
		&& go tool cover -html=temp/coverage.out -o temp/coverage.html \
		&& echo "\n\nUse to open coverage report: ${akc_cmd}firefox temp/coverage.html\n\n${akc_default}" \
		&& staticcheck -fail all -tests ./...

test-watch: temp/setup-dev
	@$(REFLEX_CMD) make test

# NOTE: Some of it is required by vscode
temp/setup-dev: temp/temp Makefile
	@go install -v golang.org/x/tools/cmd/goimports@latest \
		&& go install -v github.com/cespare/reflex@latest \
		&& go install -v golang.org/x/tools/gopls@latest \
		&& go install -v github.com/ramya-rao-a/go-outline@latest \
		&& go install -v honnef.co/go/tools/cmd/staticcheck@latest \
		&& go install -v github.com/caarlos0/svu@latest \
		&& touch temp/setup-dev

temp/temp:
	@mkdir -p temp && touch temp/temp

prepare: temp/setup-dev

clean:
	@rm -rf temp

upgrade:
	go get -u ./... && go mod tidy
