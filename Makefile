MATRIX_OS ?= darwin linux windows
MATRIX_ARCH ?= amd64 386

GIT_HASH ?= $(shell git show -s --format=%h)
APP_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

GO_DEBUG_ARGS = -v -ldflags "-X main.version=$(GO_APP_VERSION)+debug -X main.gitHash=$(GIT_HASH) -X main.buildDate=$(APP_DATE)"
GO_RELEASE_ARGS = -v -ldflags "-X main.version=$(GO_APP_VERSION) -X main.gitHash=$(GIT_HASH) -X main.buildDate=$(APP_DATE) -s -w" -tags release

MATRIX_WRAPPER ?= nodep afero
MATRIX_COMPRESSION ?= deflate gzip lzw none snappy snappystream zlib
# MATRIX_COMPRESSION ?= nocompress deflate gzip lzw snappy zlib

_TEST_FILES := $(shell find ./test -type f)
_TEST_CASES := $(patsubst %.sh,%,$(patsubst ./test-cases/%,%,$(shell find ./test-cases -type f -name '*.sh')))

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

.PHONY: install
install: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/goembed $(REQ) $(_SRC) | $(USE)
	$(eval PARTS := $(subst /, ,$*))
	$(eval BUILD := $(word 1,$(PARTS)))
	$(eval OS    := $(word 2,$(PARTS)))
	$(eval ARCH  := $(word 3,$(PARTS)))
	$(eval BIN   := $(word 4,$(PARTS)))
	$(eval ARGS  := $(if $(findstring debug,$(BUILD)),$(DEBUG_ARGS),$(RELEASE_ARGS)))

	CGO_ENABLED=$(CGO_ENABLED) GOOS="$(OS)" GOARCH="$(ARCH)" go install $(ARGS) "./cmd/..."

.PHONY: upx
upx: $(patsubst artifacts/build/%,artifacts/upx/%.upx,$(addprefix artifacts/build/release/,$(_STEMS)))

.PHONY: clean
clean::
	$(RM) -r artifacts/generated

artifacts/upx/%.upx: artifacts/build/%
	-@mkdir -p "$(@D)"
	-$(RM) -f "$(@)"
	upx -o "$@" "$<"

.PHONY: run
run: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/goembed
	"$<" $(RUN_ARGS)

.SECONDARY: $(foreach COMPRESSION,$(MATRIX_COMPRESSION),$(foreach WRAPPER,$(MATRIX_WRAPPER),artifacts/generated/compression/$(WRAPPER)/$(COMPRESSION)/compression.go))

MISSPELL := artifacts/misspell/bin/misspell
$(MISSPELL):
	@mkdir -p "$(shell pwd -P)/$(@D)"
	GO111MODULE=off GOBIN="$(shell pwd -P)/$(@D)" go get -u github.com/client9/misspell/cmd/misspell

GOMETALINTER := artifacts/gometalinter/bin/gometalinter
.PRECIOUS: $(GOMETALINTER)
$(GOMETALINTER):
	@mkdir -p "$(shell pwd -P)/$(@D)"
	GO111MODULE=off GOBIN="$(shell pwd -P)/$(@D)" go get -u github.com/alecthomas/gometalinter
	GO111MODULE=off GOBIN="$(shell pwd -P)/$(@D)" $(GOMETALINTER) --install 2>/dev/null

GOLANGCILINT := artifacts/golangci-lint/bin/golangci-lint
$(GOLANGCILINT):
	@mkdir -p "$(shell pwd -P)/$(@D)"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(shell pwd -P)/$(dir $(@D))bin" v1.21.0

.PHONY: test-compression
test-compression: $(foreach COMPRESSION,$(MATRIX_COMPRESSION),$(foreach WRAPPER,$(MATRIX_WRAPPER),artifacts/generated/compression/$(WRAPPER)/$(COMPRESSION)/test.patch))

artifacts/generated/compression/%/compression.go: artifacts/generated/compression/%/main.go $(_TEST_FILES)
	@mkdir -p "$(@D)"
	make run RUN_ARGS="./test -c "$(notdir $(*))" -w "$(subst /,,$(dir $(*)))" -f "$(@)" -p "main" -d"
	go test "$(@)"

artifacts/generated/compression/%/main.go: test/main.go.src
	@mkdir -p "$(@D)"
	cp "$(<)" "$(@)"

artifacts/generated/compression/%/a.out: artifacts/generated/compression/%/main.go artifacts/generated/compression/%/compression.go artifacts/generated/compression/%/lint
	@mkdir -p "$(@D)"
	cd "artifacts/generated/compression/$(*)" && go build -ldflags="-s -w" -o a.out .

artifacts/generated/compression/%/index.html: artifacts/generated/compression/%/a.out
	@mkdir -p "$(@D)"
	"$(@D)/a.out" | tee "$(@)"
	
artifacts/generated/compression/%/test.patch: artifacts/generated/compression/%/index.html
	@mkdir -p "$(@D)"
	diff -u "test/index.html" "$(@D)/index.html" | tee "$(@)"

artifacts/generated/compression/%/lint: $(MISSPELL) $(GOLANGCILINT)
	@mkdir -p "$(@D)"

	go vet "./$(@D)/." | tee "$@"
	! go fmt "./$(@D)/." | tee -a "$@" | grep ^

	$(MISSPELL) -w -error -locale US "./$(@D)/." | tee -a "$@"

	$(GOLANGCILINT) run ./... | tee -a "$@"

.PHONY: examples
examples: examples/webserver/assets/assets.go examples/webserver-afero/assets/assets.go

examples/webserver/assets/assets.go:
	@mkdir -p "$(@D)"
	make run RUN_ARGS="./test -c deflate -w nodep -f "$(@)" -p 'assets'"

examples/webserver-afero/assets/assets.go:
	@mkdir -p "$(@D)"
	make run RUN_ARGS="./test -c deflate -w afero -f "$(@)" -p 'assets'"

.PHONY: test-cases
test-cases: $(addprefix artifacts/test-cases/,$(_TEST_CASES))

artifacts/test-cases/%: artifacts/build/debug/$(GOOS)/$(GOARCH)/goembed
	CMD="artifacts/build/debug/$(GOOS)/$(GOARCH)/goembed" TEST_PATH="$(shell pwd)/test" bash "./test-cases/$(*).sh"
