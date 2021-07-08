#
# Build
#

BINDIR := bin
CMDDIR := cmd

SOURCES := $(shell find * -name "*.go" -or  -name "go.mod" -or -name "go.sum" \
	-or -name "Makefile")
BINS := $(shell test -d "$(CMDDIR)" && cd "$(CMDDIR)" && \
	find * -maxdepth 0 -type d -exec echo $(BINDIR)/{} \;)

.PHONY: build
build: $(BINS)

$(BINS): $(BINDIR)/%: $(SOURCES)
	mkdir -p "$(dir $@)"
	cd "$(CMDDIR)/$*" && go build -a -o "$(CURDIR)/$(BINDIR)/$*"

#
# Rules
#

rules: spammy-recruiters.json

.PHONY: spammy-recruiters.json
spammy-recruiters.json: bin/spammy-recruiters
	bin/spammy-recruiters -o "$@"
