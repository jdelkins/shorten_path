BASENAME := shorten_path
LDFLAGS := -s -w
TARGS := linux/amd64 linux/arm linux/arm64 freebsd/amd64 openbsd/amd64 netbsd/amd64 darwin/amd64 illumos/amd64
BINARIES := $(patsubst %,$(BASENAME)_%,$(subst /,_,$(TARGS)))

all: $(BINARIES)

$(BINARIES): *.go */*.go
	@$(foreach binary,$@, \
		echo building $(binary); \
		$(let os arch,$(subst _, ,$(binary:$(BASENAME)_%=%)),GOOS=$(os) GOARCH=$(arch) go build -o $(binary) -ldflags "$(LDFLAGS)") \
	)

clean:
	rm -f $(BINARIES)

.PHONY: all clean
