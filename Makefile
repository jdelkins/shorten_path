BASENAME := shorten_path
LDFLAGS := -s -w
TARGS := linux/amd64 linux/arm freebsd/amd64 openbsd/amd64 netbsd/amd64 darwin/amd64
BINARIES := $(patsubst %,$(BASENAME)_%,$(subst /,_,$(TARGS)))

all: $(BINARIES)

$(BINARIES): *.go */*.go
	gox -osarch="$(TARGS)" -ldflags="$(LDFLAGS)"
	upx $(BASENAME)_linux_amd64

clean:
	rm -f $(BINARIES)

.PHONY: all clean
