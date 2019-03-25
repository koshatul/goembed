MATRIX_OS ?= darwin linux windows
MATRIX_ARCH ?= amd64 386

GIT_HASH ?= $(shell git show -s --format=%h)
GIT_TAG ?= $(shell git tag -l --merged $(GIT_HASH) | tail -n1)
APP_VERSION ?= $(if $(TRAVIS_TAG),$(TRAVIS_TAG),$(if $(GIT_TAG),$(GIT_TAG),$(GIT_HASH)))
APP_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

MATRIX_COMPRESSION ?= nocompress deflate gzip lzw snappy zlib

_TEST_FILES := $(shell find ./test -type f)
_TEST_CASES := $(patsubst %.sh,%,$(patsubst ./test-cases/%,%,$(shell find ./test-cases -type f -name '*.sh')))

-include artifacts/make/go/Makefile

artifacts/make/%/Makefile:
	curl -sf https://jmalloc.github.io/makefiles/fetch | bash /dev/stdin $*

.PHONY: install
install: vendor $(REQ) $(_SRC) | $(USE)
	$(eval PARTS := $(subst /, ,$*))
	$(eval BUILD := $(word 1,$(PARTS)))
	$(eval OS    := $(word 2,$(PARTS)))
	$(eval ARCH  := $(word 3,$(PARTS)))
	$(eval BIN   := $(word 4,$(PARTS)))
	$(eval ARGS  := $(if $(findstring debug,$(BUILD)),$(DEBUG_ARGS),$(RELEASE_ARGS)))

	CGO_ENABLED=$(CGO_ENABLED) GOOS="$(OS)" GOARCH="$(ARCH)" go install $(ARGS) "./src/cmd/..."

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
run: artifacts/build/debug/$(GOOS)/$(GOARCH)/goembed
	$< $(RUN_ARGS)

.PHONY: test-compression
test-compression: $(addsuffix /test.patch,$(addprefix artifacts/generated/compression/,$(MATRIX_COMPRESSION)))

artifacts/generated/compression/%/compression.go: src/embed/%.go artifacts/generated/compression/%/main.go $(_TEST_FILES)
	@mkdir -p "$(@D)"
	make run RUN_ARGS="./test -c "$(*)" -f "$(@)" -p "main" -d"
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

artifacts/generated/compression/%/lint:
	@mkdir -p "$(@D)"

	go vet "./$(@D)/." | tee "$@"
	! go fmt "./$(@D)/." | tee -a "$@" | grep ^

	$(MISSPELL) -w -error -locale US "./$(@D)/." | tee -a "$@"

	$(GOMETALINTER) --disable-all --deadline=60s \
		--enable=vet \
		--enable=vetshadow \
		--enable=ineffassign \
		--enable=deadcode \
		--enable=gofmt \
		"./$(@D)/." | tee -a "$@"

	-$(GOMETALINTER) --disable-all --deadline=60s --cyclo-over=15 \
		--enable=golint \
		--enable=goconst \
		--enable=gocyclo \
		"./$(@D)/." | tee -a "$@"

.PHONY: examples
examples: examples/webserver/assets/assets.go

examples/webserver/assets/assets.go:
	@mkdir -p "$(@D)"
	make run RUN_ARGS="./test -f "$(@)" -p 'assets'"

.PHONY: test-cases
test-cases: $(addprefix artifacts/test-cases/,$(_TEST_CASES))

artifacts/test-cases/%: artifacts/build/debug/$(GOOS)/$(GOARCH)/goembed
	CMD="artifacts/build/debug/$(GOOS)/$(GOARCH)/goembed" TEST_PATH="$(shell pwd)/test" bash "./test-cases/$(*).sh"
