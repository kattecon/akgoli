REFLEX_CMD := reflex -d none -s -r "(\.go$$|^Makefile$$|^go.mod$$|^go.sum$$)" --

# Don't treat those as files
.PHONY: upgrade clean distclean build test test-watch prepare

# Disable implicit rules
.SUFFIXES:

all: build test

build: temp/setup-dev
	@go build ./...

test: temp/setup-dev
	@go test -timeout 5s -count=1 -shuffle on -coverprofile=temp/coverage.out \
			-ldflags="-X \"github.com/kattecon/akgoli/appinfo.version=3.2.1\" -X \"github.com/kattecon/akgoli/appinfo.idName=test\"" \
			-cover ./... \
		&& echo 'Running: go tool cover' \
		&& go tool cover -html=temp/coverage.out -o temp/coverage.html \
		&& echo "\n\nUse to open coverage report: ${akc_cmd}firefox temp/coverage.html\n\n${akc_default}" \
		&& echo "Running: staticcheck" \
		&& go tool staticcheck -fail all -tests ./... \
		&& echo "Running: go vet" \
		&& go vet ./...

test-watch: temp/setup-dev
	@$(REFLEX_CMD) make test

temp/setup-dev: temp/temp Makefile
	@touch temp/setup-dev

temp/temp:
	@mkdir -p temp && touch temp/temp

prepare: temp/setup-dev

clean:
	@rm -rf temp

upgrade:
	go get -u ./... && go mod tidy
